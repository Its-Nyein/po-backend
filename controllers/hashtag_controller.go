package controllers

import (
	"net/http"
	"strconv"

	"po-backend/services"

	"github.com/labstack/echo/v5"
)

type HashtagController struct {
	Service *services.HashtagService
}

func NewHashtagController(service *services.HashtagService) *HashtagController {
	return &HashtagController{Service: service}
}

func (ctrl *HashtagController) GetPostsByHashtag(c *echo.Context) error {
	tag := c.Param("tag")
	if tag == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Tag is required"})
	}

	posts, err := ctrl.Service.GetPostsByHashtag(tag)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, posts)
}

func (ctrl *HashtagController) GetTrending(c *echo.Context) error {
	limitStr := c.QueryParam("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	trending, err := ctrl.Service.GetTrending(limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, trending)
}
