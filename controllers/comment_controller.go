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
		if _, err := ctrl.NotiService.Create("comment", content, userID, req.PostID); err != nil {
			log.Printf("Failed to create notification: %v", err)
		}
		utilities.SendWebSocketMessage(post.UserID, "notis")
	}

	return c.JSON(http.StatusCreated, comment)
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
