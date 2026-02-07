package controllers

import (
	"net/http"
	"strconv"

	"po-backend/services"

	"github.com/labstack/echo/v5"
)

type NotificationController struct {
	Service *services.NotificationService
}

func NewNotificationController(service *services.NotificationService) *NotificationController {
	return &NotificationController{Service: service}
}

func (ctrl *NotificationController) GetNotifications(c *echo.Context) error {
	userID := c.Get("userID").(uint)
	notifications, err := ctrl.Service.GetByPostOwner(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, notifications)
}

func (ctrl *NotificationController) MarkAllRead(c *echo.Context) error {
	userID := c.Get("userID").(uint)
	if err := ctrl.Service.MarkAllRead(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "All notifications marked as read"})
}

func (ctrl *NotificationController) MarkOneRead(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid notification ID"})
	}

	if err := ctrl.Service.MarkOneRead(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Notification marked as read"})
}
