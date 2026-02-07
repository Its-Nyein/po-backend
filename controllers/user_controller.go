package controllers

import (
	"net/http"
	"strconv"

	"po-backend/dtos"
	"po-backend/services"
	"po-backend/validation"

	"github.com/labstack/echo/v5"
)

type UserController struct {
	Service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{Service: service}
}

func (ctrl *UserController) GetUsers(c *echo.Context) error {
	users, err := ctrl.Service.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

func (ctrl *UserController) GetUserByID(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	user, err := ctrl.Service.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) RegisterUser(c *echo.Context) error {
	var req dtos.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := validation.ValidateStruct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "Validation failed",
			"fields": validation.FormatValidationErrors(err),
		})
	}

	user, err := ctrl.Service.Register(req.Name, req.Username, req.Bio, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username already exists"})
	}

	return c.JSON(http.StatusCreated, user)
}

func (ctrl *UserController) LoginUser(c *echo.Context) error {
	var req dtos.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := validation.ValidateStruct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "Validation failed",
			"fields": validation.FormatValidationErrors(err),
		})
	}

	user, token, err := ctrl.Service.Login(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, dtos.LoginResponse{
		Token: token,
		User:  user,
	})
}

func (ctrl *UserController) VerifyToken(c *echo.Context) error {
	userID := c.Get("userID").(uint)
	user, err := ctrl.Service.GetByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	return c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) SearchUsers(c *echo.Context) error {
	query := c.QueryParam("q")
	users, err := ctrl.Service.Search(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}
