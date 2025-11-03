package http

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/root-N-root/webipfs/types"
)

// TODO: chan
func Run(fuCh chan types.FileUpdate, msgCh chan string) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ready")
	})
	log.Fatal(app.Listen(fmt.Sprintf(":%d", types.PORT)))
}
