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

type StoryController struct {
	Service       *services.StoryService
	FollowService *services.FollowService
}

func NewStoryController(service *services.StoryService, followService *services.FollowService) *StoryController {
	return &StoryController{
		Service:       service,
		FollowService: followService,
	}
}

func (ctrl *StoryController) CreateStory(c *echo.Context) error {
	var req dtos.CreateStoryRequest
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
	story, err := ctrl.Service.Create(req.Content, req.Privacy, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Broadcast to followers
	followerIDs, _ := ctrl.FollowService.GetFollowerIDs(userID)
	for _, fid := range followerIDs {
		utilities.SendWebSocketMessage(fid, "stories")
	}

	return c.JSON(http.StatusCreated, story)
}

func (ctrl *StoryController) DeleteStory(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid story ID"})
	}

	if err := ctrl.Service.Delete(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Story deleted"})
}

func (ctrl *StoryController) GetFeedStories(c *echo.Context) error {
	userID := c.Get("userID").(uint)
	stories, err := ctrl.Service.GetFeedStories(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, stories)
}

func (ctrl *StoryController) GetUserStories(c *echo.Context) error {
	targetUserID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	viewerID := c.Get("userID").(uint)
	stories, err := ctrl.Service.GetUserStories(uint(targetUserID), viewerID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, stories)
}

func (ctrl *StoryController) RecordView(c *echo.Context) error {
	storyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid story ID"})
	}

	viewerID := c.Get("userID").(uint)
	view, err := ctrl.Service.CreateView(uint(storyID), viewerID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, view)
}

func (ctrl *StoryController) GetViewers(c *echo.Context) error {
	storyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid story ID"})
	}

	views, err := ctrl.Service.GetViewers(uint(storyID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, views)
}
