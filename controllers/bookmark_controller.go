package controllers

import (
	"net/http"
	"strconv"

	"po-backend/services"

	"github.com/labstack/echo/v5"
)

type BookmarkController struct {
	Service *services.BookmarkService
}

func NewBookmarkController(service *services.BookmarkService) *BookmarkController {
	return &BookmarkController{Service: service}
}

func (ctrl *BookmarkController) CreateBookmark(c *echo.Context) error {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
	}

	userID := c.Get("userID").(uint)
	bookmark, err := ctrl.Service.Create(uint(postID), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, bookmark)
}

func (ctrl *BookmarkController) DeleteBookmark(c *echo.Context) error {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
	}

	userID := c.Get("userID").(uint)
	if err := ctrl.Service.Delete(uint(postID), userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Bookmark removed"})
}

func (ctrl *BookmarkController) GetBookmarks(c *echo.Context) error {
	userID := c.Get("userID").(uint)
	bookmarks, err := ctrl.Service.GetByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, bookmarks)
}
