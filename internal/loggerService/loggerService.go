package loggerService

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError // Default to 500

	// Determine if the error is an HTTP error
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	// Log the error
	log.Printf("Error occurred: %v", err)

	// Handle only 500 internal server errors or unexpected errors
	if code == http.StatusInternalServerError {
		// Send messages to Telegram and email only for 500 errors
		go sendMessageToTelegramService(err, c)
		// go sendEmail(err, c) // todo - implement email sending
	}

	// Send the error response to the client
	if code == http.StatusInternalServerError {
		c.JSON(code, echo.Map{"message": "Internal Server Error"})
	} else {
		// For all other errors, use the standard Echo error handling
		c.JSON(code, echo.Map{"message": err.Error()})
	}
}

// Helper function to format request headers for logging
func formatRequestHeaders(headers http.Header) string {
	var builder strings.Builder
	for name, values := range headers {
		// Skip sensitive headers
		if strings.ToLower(name) == "authorization" {
			continue
		}
		// Combine multiple header values with a comma
		valueString := strings.Join(values, ", ")
		fmt.Fprintf(&builder, "%s: %s\n", name, valueString)
	}
	return builder.String()
}

func sendMessageToTelegramService(err error, c echo.Context) {
	applicationName := os.Getenv("APPLICATION_NAME")
	botToken := os.Getenv("TELEGRAM_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	// Gather additional information
	requestURL := c.Request().URL.String()
	method := c.Request().Method
	headers := formatRequestHeaders(c.Request().Header)

	// Create a detailed error message
	text := fmt.Sprintf(
		"%s\n\nError occurred: %v\nMethod: %s\nURL: %s\nHeaders: %s",
		applicationName, err, method, requestURL, headers,
	)
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", botToken, chatID, url.QueryEscape(text))

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Printf("Failed to send message to Telegram: %v", err)
		return
	}
	defer resp.Body.Close()

	// Read the body of the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read Telegram response body: %v", err)
		return
	}

	// Check if the response status code is not 200 OK
	if resp.StatusCode != http.StatusOK {
		log.Printf("Telegram API responded with non-OK status: %s, body: %s", resp.Status, string(body))
		return
	}

	// Optionally, log a success message or do further processing with 'body'
	log.Println("Successfully sent message to Telegram")
}
