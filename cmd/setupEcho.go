package main

import (
	"github.com/go-playground/validator"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
	"todo-list-service/internal/loggerService"
	mdv "todo-list-service/internal/mdw"
)

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

	// Custom validator for request body
	e.Validator = &CustomValidator{validator: validator.New()}

	// Custom HTTP error handler
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
			unprotectedEndpoints := []string{
				"/health",
				"/",
				"/authentication/login",
				"/authentication/register",
			}

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

	//origins := strings.Split(allowedOrigins, ",")
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		//AllowOrigins: origins,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch},
		AllowHeaders: []string{
			echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,
			echo.HeaderAuthorization, "X-Requested-With",
			//echo.MIMEMultipartForm, // frontend is not working with this rule
		},
		AllowCredentials: true, // Set to true if your frontend sends credentials like cookies or auth headers
	}))

	return e
}
