package middlewares

import (
	"net/http"
	"strconv"

	"po-backend/services"

	"github.com/labstack/echo/v5"
)

func IsPostOwner(postService *services.PostService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			userID := c.Get("userID").(uint)
			postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Invalid post ID",
				})
			}

			if !postService.IsOwner(uint(postID), userID) {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "You are not the owner of this post",
				})
			}

			return next(c)
		}
	}
}

func IsStoryOwner(storyService *services.StoryService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			userID := c.Get("userID").(uint)
			storyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Invalid story ID",
				})
			}

			if !storyService.IsOwner(uint(storyID), userID) {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "You are not the owner of this story",
				})
			}

			return next(c)
		}
	}
}

func IsCommentOwner(commentService *services.CommentService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			userID := c.Get("userID").(uint)
			commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Invalid comment ID",
				})
			}

			if !commentService.IsOwner(uint(commentID), userID) {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "You are not the owner of this comment",
				})
			}

			return next(c)
		}
	}
}
