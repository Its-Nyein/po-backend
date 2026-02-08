package main

import (
	"log"
	"net/http"

	"po-backend/configs"
	"po-backend/routes"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	cfg := configs.Envs

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{cfg.CORSOrigin},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	if err := cfg.ConnectDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := cfg.ConnectRedis(); err != nil {
		log.Println("Warning: Failed to connect to Redis:", err)
	}

	if err := cfg.InitializeDB(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Po API is running"})
	})

	routes.InitializeRoutes(e, cfg.DB)

	log.Println("Starting server on :" + cfg.ServerPort)
	if err := e.Start(":" + cfg.ServerPort); err != nil && err != http.ErrServerClosed {
		log.Fatal("Failed to start server:", err)
	}
}
