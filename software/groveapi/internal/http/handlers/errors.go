package handlers

import "github.com/gofiber/fiber/v2"

type APIError struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

func SendError(c *fiber.Ctx, status int, code string, message string) error {
	return c.Status(status).JSON(APIError{
		Error:   code,
		Message: message,
	})
}

func SendErrorDetail(c *fiber.Ctx, status int, code string, message string, detail string) error {
	return c.Status(status).JSON(APIError{
		Error:   code,
		Message: message,
		Detail:  detail,
	})
}
