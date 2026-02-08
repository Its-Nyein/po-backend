package controllers

import (
	"log"
	"net/http"
	"strconv"

	"po-backend/dtos"
	"po-backend/services"
	"po-backend/utilities"
	"po-backend/validation"

	"github.com/labstack/echo/v5"
)

type PostController struct {
	Service        *services.PostService
	FollowService  *services.FollowService
	HashtagService *services.HashtagService
	NotiService    *services.NotificationService
	UserService    *services.UserService
}

func NewPostController(
	service *services.PostService,
	followService *services.FollowService,
	hashtagService *services.HashtagService,
	notiService *services.NotificationService,
	userService *services.UserService,
) *PostController {
	return &PostController{
		Service:        service,
		FollowService:  followService,
		HashtagService: hashtagService,
		NotiService:    notiService,
		UserService:    userService,
	}
}

func (ctrl *PostController) GetAllPosts(c *echo.Context) error {
	posts, err := ctrl.Service.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, posts)
}

func (ctrl *PostController) GetPostByID(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
	}

	post, err := ctrl.Service.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Post not found"})
	}

	return c.JSON(http.StatusOK, post)
}

func (ctrl *PostController) CreatePost(c *echo.Context) error {
	var req dtos.CreatePostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := validation.ValidateStruct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "Validation failed",
			"fields": validation.FormatValidationErrors(err),
		})
	}

	userID := c.Get("userID").(uint)
	post, err := ctrl.Service.Create(req.Content, userID, req.QuotedPostID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if err := ctrl.HashtagService.SyncPostHashtags(post.ID, req.Content); err != nil {
		log.Printf("Failed to sync hashtags: %v", err)
	}

	ctrl.processMentions(req.Content, userID, post.ID)

	if req.QuotedPostID != nil {
		original, _ := ctrl.Service.GetByID(*req.QuotedPostID)
		if original != nil && original.UserID != userID {
			actor, _ := ctrl.UserService.GetByID(userID)
			notiContent := "quoted your post"
			if actor != nil {
				notiContent = actor.Name + " quoted your post"
			}
			pid := *req.QuotedPostID
			if _, err := ctrl.NotiService.Create("quote", notiContent, userID, original.UserID, &pid); err != nil {
				log.Printf("Failed to create quote notification: %v", err)
			}
			utilities.SendWebSocketMessage(original.UserID, "notis")
		}
	}

	return c.JSON(http.StatusCreated, post)
}

func (ctrl *PostController) UpdatePost(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
	}

	var req dtos.UpdatePostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := validation.ValidateStruct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "Validation failed",
			"fields": validation.FormatValidationErrors(err),
		})
	}

	post, err := ctrl.Service.Update(uint(id), req.Content)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if err := ctrl.HashtagService.SyncPostHashtags(uint(id), req.Content); err != nil {
		log.Printf("Failed to sync hashtags: %v", err)
	}

	return c.JSON(http.StatusOK, post)
}

func (ctrl *PostController) DeletePost(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
	}

	if err := ctrl.Service.Delete(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Post deleted"})
}

func (ctrl *PostController) GetFollowingPosts(c *echo.Context) error {
	userID := c.Get("userID").(uint)
	followingIDs, err := ctrl.FollowService.GetFollowingIDs(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if len(followingIDs) == 0 {
		return c.JSON(http.StatusOK, []interface{}{})
	}

	posts, err := ctrl.Service.GetByUserIDs(followingIDs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, posts)
}

func (ctrl *PostController) processMentions(content string, actorID uint, postID uint) {
	usernames := utilities.ParseMentions(content)
	actor, _ := ctrl.UserService.GetByID(actorID)

	for _, username := range usernames {
		mentioned, err := ctrl.UserService.GetByUsername(username)
		if err != nil || mentioned.ID == actorID {
			continue
		}
		notiContent := "mentioned you in a post"
		if actor != nil {
			notiContent = actor.Name + " mentioned you in a post"
		}
		pid := postID
		if _, err := ctrl.NotiService.Create("mention", notiContent, actorID, mentioned.ID, &pid); err != nil {
			log.Printf("Failed to create mention notification: %v", err)
		}
		utilities.SendWebSocketMessage(mentioned.ID, "notis")
	}
}
