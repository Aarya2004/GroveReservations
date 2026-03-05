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

type AuthHTTP struct {
	DB *gorm.DB
	SB *supabase.Client
}

func NewAuthHTTP(db *gorm.DB, sb *supabase.Client) *AuthHTTP { return &AuthHTTP{DB: db, SB: sb} }

type AdminCreateUserInput struct {
	Email       string `json:"email"`
	Password    string `json:"password,omitempty"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	VillaNumber int    `json:"villa_number"`
	PhoneNumber string `json:"phone_number"`
}

func (h *AuthHTTP) AdminCreateUser(c *fiber.Ctx) error {
	var in AdminCreateUserInput
	if err := c.BodyParser(&in); err != nil {
		return SendError(c, http.StatusBadRequest, "invalid_json", "could not parse request body")
	}
	if in.Email == "" || in.Name == "" {
		return SendError(c, http.StatusBadRequest, "missing_required_fields", "email and name are required")
	}

	params := types.AdminCreateUserRequest{
		Email:        in.Email,
		UserMetadata: map[string]any{"name": in.Name},
		EmailConfirm: false,
	}
	if in.Password != "" {
		params.Password = &in.Password
	}

	u, err := h.SB.Auth.AdminCreateUser(params)
	if err != nil {
		return SendErrorDetail(c, http.StatusInternalServerError, "auth_admin_create_failed", "failed to create auth user", err.Error())
	}
	authID := u.ID.String()
	authEmail := u.Email
	uid, err := uuid.Parse(authID)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, "auth_id_parse_failed", "could not parse auth user ID")
	}

	profile := store.User{
		ID:          uid,
		Email:       authEmail,
		Name:        in.Name,
		Role:        in.Role,
		VillaNumber: in.VillaNumber,
		PhoneNumber: in.PhoneNumber,
	}
	if err := h.DB.WithContext(context.Background()).Create(&profile).Error; err != nil {
		return SendError(c, http.StatusInternalServerError, "profile_insert_failed", "failed to create user profile")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"id": authID, "email": authEmail})
}

type RegisterInput struct {
	Email       string `json:"email"`
	Password    string `json:"password,omitempty"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	VillaNumber int    `json:"villa_number"`
	PhoneNumber string `json:"phone_number"`
}

func (h *AuthHTTP) Register(c *fiber.Ctx) error {
	var in RegisterInput
	if err := c.BodyParser(&in); err != nil {
		return SendError(c, http.StatusBadRequest, "invalid_json", "could not parse request body")
	}
	if in.Email == "" || in.Name == "" {
		return SendError(c, http.StatusBadRequest, "missing_required_fields", "email and name are required")
	}

	params := types.SignupRequest{
		Email: in.Email,
	}
	if in.Password != "" {
		params.Password = in.Password
	}

	u, err := h.SB.Auth.Signup(params)
	if err != nil {
		return SendErrorDetail(c, http.StatusInternalServerError, "user_registration_failed", "failed to register user", err.Error())
	}
	authID := u.ID.String()
	authEmail := u.Email
	uid, err := uuid.Parse(authID)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, "auth_id_parse_failed", "could not parse auth user ID")
	}

	profile := store.User{
		ID:          uid,
		Email:       authEmail,
		Name:        in.Name,
		Role:        in.Role,
		VillaNumber: in.VillaNumber,
		PhoneNumber: in.PhoneNumber,
	}
	if err := h.DB.WithContext(context.Background()).Create(&profile).Error; err != nil {
		return SendError(c, http.StatusInternalServerError, "profile_insert_failed", "failed to create user profile")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"id": authID, "email": authEmail})
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHTTP) Login(c *fiber.Ctx) error {
	var in LoginInput
	if err := c.BodyParser(&in); err != nil {
		return SendError(c, http.StatusBadRequest, "invalid_json", "could not parse request body")
	}
	if in.Email == "" || in.Password == "" {
		return SendError(c, http.StatusBadRequest, "missing_required_fields", "email and password are required")
	}

	tokenResponse, err := h.SB.Auth.SignInWithEmailPassword(in.Email, in.Password)
	if err != nil {
		return SendErrorDetail(c, http.StatusBadRequest, "invalid_credentials", "invalid email or password", err.Error())
	}
	sess := tokenResponse.Session
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
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"token_type":    "bearer",
		"access_token":  sess.AccessToken,
		"refresh_token": sess.RefreshToken,
		"user": fiber.Map{
			"id":    usr.ID.String(),
			"email": usr.Email,
		},
	})
}

func (h *AuthHTTP) Logout(c *fiber.Ctx) error {
	err := h.SB.Auth.Logout()
	if err != nil {
		return SendErrorDetail(c, http.StatusInternalServerError, "error_logging_user_out", "failed to log out", err.Error())
	}
	return c.SendStatus(http.StatusNoContent)
}

func (h *AuthHTTP) GetCurrentUser(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	var profile store.User
	if err := h.DB.Where("id = ?", userID).First(&profile).Error; err != nil {
		return SendError(c, http.StatusNotFound, "user_not_found", "user profile not found")
	}
	return c.Status(http.StatusOK).JSON(profile)
}

func (h *AuthHTTP) ListUsers(c *fiber.Ctx) error {
	var users []store.User
	result := h.DB.Find(&users)
	if result.Error != nil {
		return SendErrorDetail(c, http.StatusInternalServerError, "error_listing_users", "failed to list users", result.Error.Error())
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
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		return nil
	}

	var in UpdateUserInput
	if err := c.BodyParser(&in); err != nil {
		return SendError(c, http.StatusBadRequest, "invalid_json", "could not parse request body")
	}

	var user store.User
	if err := h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return SendError(c, http.StatusNotFound, "user_not_found", "user not found")
	}

	result := h.DB.Model(&user).Updates(store.User{
		Name: in.Name, Role: in.Role, VillaNumber: in.VillaNumber, PhoneNumber: in.PhoneNumber,
	})
	if result.Error != nil {
		return SendError(c, http.StatusInternalServerError, "update_failed", "failed to update user")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"id": user.ID})
}

func (h *AuthHTTP) DeactivateUser(c *fiber.Ctx) error {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		return nil
	}

	var user store.User
	if err := h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return SendError(c, http.StatusNotFound, "user_not_found", "user not found")
	}

	result := h.DB.Model(&user).Update("active", false)
	if result.Error != nil {
		return SendError(c, http.StatusInternalServerError, "deactivate_failed", "failed to deactivate user")
	}

	return c.SendStatus(http.StatusNoContent)
}
