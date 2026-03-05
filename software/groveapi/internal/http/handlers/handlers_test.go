package handlers_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"groveapi/internal/http/handlers"

	"github.com/gofiber/fiber/v2"
)

func setupTestApp() *fiber.App {
	app := fiber.New()
	// Resource handler without DB — will error, but we can test UUID parsing
	h := handlers.NewResourceHTTP(nil)
	app.Get("/resources/:id", h.GetResource)
	return app
}

func TestInvalidUUID_Returns400(t *testing.T) {
	app := setupTestApp()

	req, _ := http.NewRequest("GET", "/resources/not-a-uuid", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("expected 400, got %d: %s", resp.StatusCode, string(body))
	}
}

func TestInvalidUUID_ResponseContainsError(t *testing.T) {
	app := setupTestApp()

	req, _ := http.NewRequest("GET", "/resources/invalid", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "invalid_uuid") {
		t.Errorf("expected body to contain 'invalid_uuid', got: %s", string(body))
	}
}
