package controllers

import (
	"net/http"
	"strconv"

	"po-backend/dtos"
	"po-backend/services"
	"po-backend/utilities"
	"po-backend/validation"

	"github.com/labstack/echo/v5"
)

type ConversationController struct {
	Service *services.ConversationService
}

func NewConversationController(service *services.ConversationService) *ConversationController {
	return &ConversationController{Service: service}
}

func (ctrl *ConversationController) GetConversations(c *echo.Context) error {
	userID := c.Get("userID").(uint)
	conversations, err := ctrl.Service.GetConversations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, conversations)
}

func (ctrl *ConversationController) CreateConversation(c *echo.Context) error {
	var req dtos.CreateConversationRequest
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
	conv, err := ctrl.Service.GetOrCreateConversation(userID, req.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, conv)
}

func (ctrl *ConversationController) GetUnreadCount(c *echo.Context) error {
	userID := c.Get("userID").(uint)
	count, err := ctrl.Service.GetUnreadTotal(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]int64{"count": count})
}

func (ctrl *ConversationController) CheckMutualFollow(c *echo.Context) error {
	otherUserID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	userID := c.Get("userID").(uint)
	mutual, err := ctrl.Service.AreMutualFollowers(userID, uint(otherUserID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]bool{"canMessage": mutual})
}

func (ctrl *ConversationController) GetMessages(c *echo.Context) error {
	convID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid conversation ID"})
	}

	var cursor uint
	if cursorStr := c.QueryParam("cursor"); cursorStr != "" {
		cursorVal, err := strconv.ParseUint(cursorStr, 10, 32)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid cursor"})
		}
		cursor = uint(cursorVal)
	}

	userID := c.Get("userID").(uint)
	messages, err := ctrl.Service.GetMessages(uint(convID), userID, cursor)
	if err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, messages)
}

func (ctrl *ConversationController) SendMessage(c *echo.Context) error {
	convID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid conversation ID"})
	}

	var req dtos.SendMessageRequest
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
	msg, err := ctrl.Service.SendMessage(uint(convID), userID, req.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Push WebSocket events to recipient
	otherUserID, wsErr := ctrl.Service.GetOtherParticipantID(uint(convID), userID)
	if wsErr == nil {
		utilities.SendWebSocketMessage(otherUserID, "messages")
		utilities.SendWebSocketMessage(otherUserID, "conversations")
		utilities.SendWebSocketMessage(otherUserID, "unreadMessages")
	}

	return c.JSON(http.StatusCreated, msg)
}

func (ctrl *ConversationController) MarkRead(c *echo.Context) error {
	convID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid conversation ID"})
	}

	userID := c.Get("userID").(uint)
	if err := ctrl.Service.MarkRead(uint(convID), userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Marked as read"})
}
