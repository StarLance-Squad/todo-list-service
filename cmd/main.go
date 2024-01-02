package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
	"todo-list-service/internal/auth"
	"todo-list-service/internal/db"
	"todo-list-service/internal/routes"
	"todo-list-service/internal/services"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	e := setupEcho()

	// Connect to the database with GORM
	dbInstance := db.ConnectDB()
	defer func() {
		sqlDB, err := dbInstance.DB()
		if err != nil {
			log.Fatalf("Error on closing database connection: %v", err)
		}
		sqlDB.Close()
	}()

	// Initialize repositories
	userRepo := &db.GormUserRepository{DB: dbInstance}
	todoRepo := &db.GormTodoRepository{DB: dbInstance}

	// Initialize services
	svc := &services.Services{
		UserService: &services.UserService{UserRepo: userRepo},
		TodoService: &services.TodoService{TodoRepo: todoRepo},
	}

	// JWT Middleware Configuration
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set in .env file")
	}
	jwtMiddleware := middleware.JWTWithConfig(auth.JWTMiddlewareConfig(jwtSecret))

	// Initialize routes
	routes.Init(e, svc, jwtMiddleware)

	// Get server port from environment variable, default to 8000
	port := os.Getenv("DEV_PORT")
	if port == "" {
		port = "8000"
	}

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	// Start server in a goroutine
	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func setupEcho() *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS middleware
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		log.Fatal("CORS_ALLOWED_ORIGINS is not set in .env file")
	}

	origins := strings.Split(allowedOrigins, ",")
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: origins,
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	return e
}
