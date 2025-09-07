package handlers

import (
	"context"
	"groveapi/internal/store"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/supabase-community/gotrue-go/types"
	supabase "github.com/supabase-community/supabase-go"
	"gorm.io/gorm"
)

// POST `/auth/login` – Login.
// POST `/auth/logout` – Logout/invalidate session.
// GET `/users/me` – Current user profile.
// GET `/users` *(admin)* – List all users.
// PATCH `/users/:id` *(admin)* – Update role/status.
// DELETE `/users/:id` *(admin)* – Deactivate user.

type AuthHTTP struct {
	DB *gorm.DB
	SB *supabase.Client
}

func NewAuthHTTP(db *gorm.DB, sb *supabase.Client) *AuthHTTP { return &AuthHTTP{DB: db, SB: sb} }

type AdminCreateUserInput struct {
	Email       string `json:"email"`
	Password    string `json:"password,omitempty"` // optional if using magic link
	Name        string `json:"name"`
	Role        string `json:"role"`
	VillaNumber int    `json:"villa_number"`
	PhoneNumber string `json:"phone_number"`
}

func (h *AuthHTTP) AdminCreateUser(c *fiber.Ctx) error {
	var in AdminCreateUserInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid_json"})
	}
	if in.Email == "" || in.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "missing_required_fields"})
	}

	// 1) Create Auth user (Admin)
	
	params := types.AdminCreateUserRequest{
		Email:         in.Email,
		UserMetadata: map[string]any{"name": in.Name},
		EmailConfirm: false, // TODO
	}
	if in.Password != "" {
		params.Password = &in.Password
	}

	u, err := h.SB.Auth.AdminCreateUser(params)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "auth_admin_create_failed", "detail": err.Error()})
	}
	authID := u.ID.String()
	authEmail := u.Email
	uid, err := uuid.Parse(authID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "auth_id_parse_failed"})
	}

	// 2) Mirror profile row in public.users with SAME id
	profile := store.User{
		ID:          uid,
		Email:       authEmail,
		Name:        in.Name,
		Role:        in.Role,
		VillaNumber: in.VillaNumber,
		PhoneNumber: in.PhoneNumber,
	}
	if err := h.DB.WithContext(context.Background()).Create(&profile).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "profile_insert_failed"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"id": authID, "email": authEmail})
}

type RegisterInput struct {
	Email       string `json:"email"`
	Password    string `json:"password,omitempty"` // optional if using magic link
	Name        string `json:"name"`
	Role        string `json:"role"`
	VillaNumber int    `json:"villa_number"`
	PhoneNumber string `json:"phone_number"`
}

func (h *AuthHTTP) Register(c *fiber.Ctx) error { 
	var in RegisterInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid_json"})
	}
	if in.Email == "" || in.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "missing_required_fields"})
	}

	// 1) Create Auth user (Admin)
	
	params := types.SignupRequest{
		Email:         in.Email,
	}
	if in.Password != "" {
		params.Password = in.Password
	}

	u, err := h.SB.Auth.Signup(params)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "user_registration_failed", "detail": err.Error()})
	}
	authID := u.ID.String()
	authEmail := u.Email
	uid, err := uuid.Parse(authID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "auth_id_parse_failed"})
	}

	// 2) Mirror profile row in public.users with SAME id
	profile := store.User{
		ID:          uid,
		Email:       authEmail,
		Name:        in.Name,
		Role:        in.Role,
		VillaNumber: in.VillaNumber,
		PhoneNumber: in.PhoneNumber,
	}
	if err := h.DB.WithContext(context.Background()).Create(&profile).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "profile_insert_failed"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"id": authID, "email": authEmail})
}

type LoginInput struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
}

func (h *AuthHTTP) Login(c *fiber.Ctx) error { 
	var in LoginInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid_json"})
	}
	if in.Email == "" || in.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "missing_required_fields"})
	}

	tokenResponse, err := h.SB.Auth.SignInWithEmailPassword(in.Email, in.Password);
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid_credentials", "detail": err.Error()})
	}
	sess  := tokenResponse.Session
	h.SB.UpdateAuthSession(sess)

	c.Cookie(&fiber.Cookie{
		Name:     "sb-refresh-token",
		Value:    sess.RefreshToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/",
	})
	
	usr := sess.User
	userID := usr.ID.String()
	userEmail := usr.Email

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"token_type":    "bearer",
		"access_token":  sess.AccessToken,
		"refresh_token": sess.RefreshToken,
		"user": fiber.Map{
			"id":    userID,
			"email": userEmail,
		},
	})
}

func (h *AuthHTTP) Logout(c *fiber.Ctx) error { 
	err := h.SB.Auth.Logout()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error_logging_user_out", "detail": err.Error()})
	}
	return c.SendStatus(http.StatusNoContent) 
}

func (h *AuthHTTP) GetCurrentUser(c *fiber.Ctx) error { 
	u, err := h.SB.Auth.GetUser()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error_fetching_user", "detail": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"id":    u.ID.String(),
		"email": u.Email,
		"metadata": u.UserMetadata,
	}) 
}

func (h *AuthHTTP) ListUsers(c *fiber.Ctx) error {
	var users []store.User
	result := h.DB.Find(&users)
	if result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error_listing_users", "detail": result.Error.Error()})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"count": result.RowsAffected,
		"users": users,
	})
}

type UpdateUserInput struct {
	Name        string `json:"name"`
	Role        string `json:"role"`
	VillaNumber int    `json:"villa_number"`
	PhoneNumber string `json:"phone_number"`
}

func (h *AuthHTTP) UpdateUser(c *fiber.Ctx) error { 
	var user store.User;
	var in UpdateUserInput;
	if err := c.BodyParser(&in); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid_json"})
	}
	user.ID = uuid.MustParse(c.Params("id"))
	if err := h.DB.First(&user).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "user_not_found"})
	}
	h.DB.Model(&user).Updates(store.User{Name: in.Name, Role: in.Role, VillaNumber: in.VillaNumber, PhoneNumber: in.PhoneNumber})
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"id": user.ID,
	})
}

func (h *AuthHTTP) DeactivateUser(c *fiber.Ctx) error  { 
	var user store.User;
	user.ID = uuid.MustParse(c.Params("id"))
	h.DB.Model(&user).Update("active", false)
	return c.SendStatus(http.StatusNoContent)
}