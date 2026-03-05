package middleware

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	supabase "github.com/supabase-community/supabase-go"
	"gorm.io/gorm"

	"groveapi/internal/store"
)

func RequireAuth(sb *supabase.Client, db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error":   "unauthorized",
				"message": "missing or invalid Authorization header",
			})
		}
		token := strings.TrimPrefix(auth, "Bearer ")

		user, err := sb.Auth.WithToken(token).GetUser()
		if err != nil || user == nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error":   "unauthorized",
				"message": "invalid or expired token",
			})
		}

		userID := user.ID.String()

		// Look up the profile to get role
		var profile store.User
		if err := db.Where("id = ?", userID).First(&profile).Error; err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error":   "unauthorized",
				"message": "user profile not found",
			})
		}

		c.Locals("user_id", userID)
		c.Locals("role", profile.Role)
		return c.Next()
	}
}

func RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if role != "ADMIN" {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error":   "forbidden",
				"message": "admin access required",
			})
		}
		return c.Next()
	}
}
