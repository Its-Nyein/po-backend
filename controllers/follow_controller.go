package controllers

import (
	"net/http"
	"strconv"

	"po-backend/services"

	"github.com/labstack/echo/v5"
)

type FollowController struct {
	Service *services.FollowService
}

func NewFollowController(service *services.FollowService) *FollowController {
	return &FollowController{Service: service}
}

func (ctrl *FollowController) FollowUser(c *echo.Context) error {
	followingID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	followerID := c.Get("userID").(uint)
	if followerID == uint(followingID) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot follow yourself"})
	}

	follow, err := ctrl.Service.Follow(followerID, uint(followingID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, follow)
}

func (ctrl *FollowController) UnfollowUser(c *echo.Context) error {
	followingID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	followerID := c.Get("userID").(uint)
	if err := ctrl.Service.Unfollow(followerID, uint(followingID)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Unfollowed"})
}
