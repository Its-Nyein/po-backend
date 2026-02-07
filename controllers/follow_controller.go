package controllers

import (
	"log"
	"net/http"
	"strconv"

	"po-backend/services"
	"po-backend/utilities"

	"github.com/labstack/echo/v5"
)

type FollowController struct {
	Service     *services.FollowService
	NotiService *services.NotificationService
	UserService *services.UserService
}

func NewFollowController(
	service *services.FollowService,
	notiService *services.NotificationService,
	userService *services.UserService,
) *FollowController {
	return &FollowController{
		Service:     service,
		NotiService: notiService,
		UserService: userService,
	}
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

	user, _ := ctrl.UserService.GetByID(followerID)
	content := "started following you"
	if user != nil {
		content = user.Name + " started following you"
	}
	if _, err := ctrl.NotiService.Create("follow", content, followerID, uint(followingID), nil); err != nil {
		log.Printf("Failed to create follow notification: %v", err)
	}
	utilities.SendWebSocketMessage(uint(followingID), "notis")

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
