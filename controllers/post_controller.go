package controllers

import (
	"net/http"
	"strconv"

	"po-backend/dtos"
	"po-backend/services"
	"po-backend/validation"

	"github.com/labstack/echo/v5"
)

type PostController struct {
	Service       *services.PostService
	FollowService *services.FollowService
}

func NewPostController(service *services.PostService, followService *services.FollowService) *PostController {
	return &PostController{Service: service, FollowService: followService}
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
	post, err := ctrl.Service.Create(req.Content, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
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
