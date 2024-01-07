package main

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
	"todo-list-service/internal/db"
	"todo-list-service/internal/loggerService"
	mdv "todo-list-service/internal/mdw"
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

type (
	CustomValidator struct {
		validator *validator.Validate
	}
)

// Validate validates the request body. Docs: https://echo.labstack.com/docs/request#validate-data
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func setupEcho() *echo.Echo {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.HTTPErrorHandler = loggerService.CustomHTTPErrorHandler

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(mdv.PaginationMiddleware)

	// Read and decode JWT token here: https://jwt.io/
	// echo JWT Middleware Configuration. Docs: https://github.com/labstack/echo-jwt
	jwtSecret := os.Getenv("JWT_SECRET")
	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(jwtSecret),
		Skipper: func(c echo.Context) bool {
			// List of endpoints that don't require JWT authentication
			unprotectedEndpoints := []string{"/health", "/", "/authentication/login", "/authentication/register"}

			// Skip JWT middleware for endpoints in the list
			for _, endpoint := range unprotectedEndpoints {
				if c.Path() == endpoint {
					return true
				}
			}

			// Important! All other endpoints require JWT authentication
			return false
		},
	}))

	// CORS mdw
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		log.Fatal("CORS_ALLOWED_ORIGINS is not set in .env file")
	}

	origins := strings.Split(allowedOrigins, ",")
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: origins,
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.MIMEMultipartForm},
	}))

	return e
}
