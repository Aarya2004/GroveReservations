package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"net/http"
)

// POST `/auth/login` – Login.
// POST `/auth/logout` – Logout/invalidate session.
// GET `/users/me` – Current user profile.
// GET `/users` *(admin)* – List all users.
// PATCH `/users/:id` *(admin)* – Update role/status.
// DELETE `/users/:id` *(admin)* – Deactivate user.

type AuthHTTP struct {
	DB *gorm.DB
}

func NewAuthHTTP(db *gorm.DB) *AuthHTTP { return &AuthHTTP{DB: db} }

func (h *AuthHTTP) Login(c *fiber.Ctx) error   { /* TODO */ return c.SendStatus(http.StatusNotImplemented) }
func (h *AuthHTTP) Logout(c *fiber.Ctx) error     { /* TODO */ return c.SendStatus(http.StatusNotImplemented) }
func (h *AuthHTTP) GetCurrentUser(c *fiber.Ctx) error     { /* TODO */ return c.SendStatus(http.StatusNotImplemented) }
func (h *AuthHTTP) ListUsers(c *fiber.Ctx) error     { /* TODO */ return c.SendStatus(http.StatusNotImplemented) }
func (h *AuthHTTP) UpdateUser(c *fiber.Ctx) error     { /* TODO */ return c.SendStatus(http.StatusNotImplemented) }
func (h *AuthHTTP) DeactivateUser(c *fiber.Ctx) error     { /* TODO */ return c.SendStatus(http.StatusNotImplemented) }