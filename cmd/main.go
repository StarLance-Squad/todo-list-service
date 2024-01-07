package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
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
	// Setup Echo
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

	// Initialize app routes
	routes.Init(e, svc)

	// Get server port from environment variable, default to 8000
	port := os.Getenv("DEV_PORT")
	if port == "" {
		port = "8000"
	}

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	e.GET("/", func(c echo.Context) error { return c.NoContent(http.StatusOK) })

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
