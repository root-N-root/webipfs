package http

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestHandlerWebsocket(t *testing.T) {
	app := fiber.New()
	app.Use("/ws", handlerWebsocket)
	app.Get("/ws", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	t.Run("should upgrade to ws", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/ws", nil)
		req.Header.Set("Connection", "Upgrade")
		req.Header.Set("Upgrade", "websocket")
		req.Header.Set("Sec-Websocket-Version", "13")
		req.Header.Set("Sec-Websocket-Key", "test-key")

		_, err := app.Test(req)
		assert.NoError(t, err)
	})
	t.Run("should reject non-ws", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/ws", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUpgradeRequired, resp.StatusCode)
	})
}
