package controllers

import (
	"log"
	"net/http"
	"strconv"

	"po-backend/services"
	"po-backend/utilities"

	"github.com/labstack/echo/v5"
)

type LikeController struct {
	Service        *services.LikeService
	PostService    *services.PostService
	CommentService *services.CommentService
	NotiService    *services.NotificationService
	UserService    *services.UserService
}

func NewLikeController(
	service *services.LikeService,
	postService *services.PostService,
	commentService *services.CommentService,
	notiService *services.NotificationService,
	userService *services.UserService,
) *LikeController {
	return &LikeController{
		Service:        service,
		PostService:    postService,
		CommentService: commentService,
		NotiService:    notiService,
		UserService:    userService,
	}
}

func (ctrl *LikeController) LikePost(c *echo.Context) error {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
	}

	userID := c.Get("userID").(uint)
	like, err := ctrl.Service.LikePost(uint(postID), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	post, _ := ctrl.PostService.GetByID(uint(postID))
	if post != nil && post.UserID != userID {
		user, _ := ctrl.UserService.GetByID(userID)
		content := "liked your post"
		if user != nil {
			content = user.Name + " liked your post"
		}
		pid := uint(postID)
		if _, err := ctrl.NotiService.Create("like", content, userID, post.UserID, &pid); err != nil {
			log.Printf("Failed to create notification: %v", err)
		}
		utilities.SendWebSocketMessage(post.UserID, "notis")
	}

	return c.JSON(http.StatusOK, like)
}

func (ctrl *LikeController) UnlikePost(c *echo.Context) error {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
	}

	userID := c.Get("userID").(uint)
	if err := ctrl.Service.UnlikePost(uint(postID), userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Post unliked"})
}

func (ctrl *LikeController) LikeComment(c *echo.Context) error {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid comment ID"})
	}

	userID := c.Get("userID").(uint)
	like, err := ctrl.Service.LikeComment(uint(commentID), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	comment, _ := ctrl.CommentService.GetByID(uint(commentID))
	if comment != nil {
		post, _ := ctrl.PostService.GetByID(comment.PostID)
		if post != nil && post.UserID != userID {
			user, _ := ctrl.UserService.GetByID(userID)
			content := "liked a comment on your post"
			if user != nil {
				content = user.Name + " liked a comment on your post"
			}
			pid := comment.PostID
			if _, err := ctrl.NotiService.Create("like", content, userID, post.UserID, &pid); err != nil {
				log.Printf("Failed to create notification: %v", err)
			}
			utilities.SendWebSocketMessage(post.UserID, "notis")
		}
	}

	return c.JSON(http.StatusOK, like)
}

func (ctrl *LikeController) UnlikeComment(c *echo.Context) error {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid comment ID"})
	}

	userID := c.Get("userID").(uint)
	if err := ctrl.Service.UnlikeComment(uint(commentID), userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Comment unliked"})
}

func (ctrl *LikeController) GetPostLikers(c *echo.Context) error {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
	}

	likes, err := ctrl.Service.GetPostLikers(uint(postID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, likes)
}

func (ctrl *LikeController) GetCommentLikers(c *echo.Context) error {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid comment ID"})
	}

	likes, err := ctrl.Service.GetCommentLikers(uint(commentID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, likes)
}
