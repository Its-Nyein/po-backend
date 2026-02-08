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

type CommentController struct {
	Service     *services.CommentService
	PostService *services.PostService
	NotiService *services.NotificationService
	UserService *services.UserService
}

func NewCommentController(
	service *services.CommentService,
	postService *services.PostService,
	notiService *services.NotificationService,
	userService *services.UserService,
) *CommentController {
	return &CommentController{
		Service:     service,
		PostService: postService,
		NotiService: notiService,
		UserService: userService,
	}
}

func (ctrl *CommentController) CreateComment(c *echo.Context) error {
	var req dtos.CreateCommentRequest
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
	comment, err := ctrl.Service.Create(req.Content, userID, req.PostID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	post, _ := ctrl.PostService.GetByID(req.PostID)
	if post != nil && post.UserID != userID {
		user, _ := ctrl.UserService.GetByID(userID)
		content := "commented on your post"
		if user != nil {
			content = user.Name + " commented on your post"
		}
		pid := req.PostID
		if _, err := ctrl.NotiService.Create("comment", content, userID, post.UserID, &pid); err != nil {
			log.Printf("Failed to create notification: %v", err)
		}
		utilities.SendWebSocketMessage(post.UserID, "notis")
	}

	ctrl.processMentions(req.Content, userID, req.PostID)

	return c.JSON(http.StatusCreated, comment)
}

func (ctrl *CommentController) processMentions(content string, actorID uint, postID uint) {
	usernames := utilities.ParseMentions(content)
	actor, _ := ctrl.UserService.GetByID(actorID)

	for _, username := range usernames {
		mentioned, err := ctrl.UserService.GetByUsername(username)
		if err != nil || mentioned.ID == actorID {
			continue
		}
		notiContent := "mentioned you in a comment"
		if actor != nil {
			notiContent = actor.Name + " mentioned you in a comment"
		}
		pid := postID
		if _, err := ctrl.NotiService.Create("mention", notiContent, actorID, mentioned.ID, &pid); err != nil {
			log.Printf("Failed to create mention notification: %v", err)
		}
		utilities.SendWebSocketMessage(mentioned.ID, "notis")
	}
}

func (ctrl *CommentController) UpdateComment(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid comment ID"})
	}

	var req dtos.UpdateCommentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := validation.ValidateStruct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "Validation failed",
			"fields": validation.FormatValidationErrors(err),
		})
	}

	comment, err := ctrl.Service.Update(uint(id), req.Content)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, comment)
}

func (ctrl *CommentController) DeleteComment(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid comment ID"})
	}

	if err := ctrl.Service.Delete(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Comment deleted"})
}
