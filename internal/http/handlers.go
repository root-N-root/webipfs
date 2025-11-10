package http

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/ipfs/go-cid"

	"github.com/root-N-root/webipfs/types"
)

// TODO: chan
func Run(ctx context.Context, con *types.Connector, client types.FileService) {
	app := fiber.New()
	//GRACEFULL shutdown
	go func() {
		<-ctx.Done()

		log.Println("Shutdown signal received, stopping server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		done := make(chan struct{})
		go func() {
			app.Shutdown()
			close(done)
		}()
		select {
		case <-done:
			log.Println("Server shut down gracefully")
		case <-shutdownCtx.Done():
			log.Println("Shutdown timeout reached, forcing exit")
		}
	}()
	app.Post("/upload", appHandlerUpload(con, client))
	app.Post("/download", appHandlerDownload(con, client))
	app.Use("/ws", handlerWebsocket)
	app.Get("/ws", websocket.New(websocketConn(con)))
	app.Static("/", "./public")
	log.Fatal(app.Listen(fmt.Sprintf(":%d", types.PORT)))
}

type DownloadRequest struct {
	Link string `json:"link"`
}

func appHandlerDownload(con *types.Connector, client types.FileService) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var req DownloadRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).SendString("Invalid request body: " + err.Error())
		}

		// Извлекаем CID из IPFS-ссылки или используем напрямую, если это просто CID
		cidStr := req.Link
		if strings.HasPrefix(req.Link, "ipfs://") {
			cidStr = strings.TrimPrefix(req.Link, "ipfs://")
			// Удаляем возможный слэш в начале
			cidStr = strings.TrimPrefix(cidStr, "/")
		}

		// Проверяем, что строка - это валидный CID
		if _, err := cid.Decode(cidStr); err != nil {
			return c.Status(400).SendString("Invalid CID: " + err.Error())
		}

		// Вызываем метод IPFS-клиента для получения файла по CID
		fileUpdate, err := client.GetFile(c.Context(), cidStr)
		if err != nil {
			return c.Status(500).SendString("Error getting file from IPFS: " + err.Error())
		}

		// Отправляем информацию о файле через соединение
		con.SendFileUp(fileUpdate)

		return c.Status(200).JSON(fiber.Map{
			"message": "Download started",
			"file":    fileUpdate,
		})
	}
}

func appHandlerUpload(con *types.Connector, client types.FileService) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}
		filePath := fmt.Sprintf("%s/%d", types.FILES_DIR, time.Now().UnixNano())
		err = c.SaveFile(file, filePath)
		if err != nil {
			return err
		}

		// con.FileChan <- types.File{Path: filePath, Name: file.Filename}
		fu, err := client.AddFile(filePath, file.Filename)
		if err != nil {
			return err
		}
		con.SendFileUp(fu)

		return nil
	}
}

func handlerWebsocket(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func websocketConn(con *types.Connector) func(c *websocket.Conn) {
	return func(c *websocket.Conn) {
		log.Println("WebSocket connection established")
		defer c.Close()

		go func() {
			defer c.Close()
			for {
				fu := <-con.FileUpHttpChan
				if err := c.WriteJSON(map[string]any{
					"type":     "file_update",
					"cid":      fu.CID,
					"name":     fu.Name,
					"peers":    fu.Peers,
					"progress": fu.Progress,
					"status":   fu.Status,
				}); err != nil {
					log.Println("write json:", err)
					return
				}
			}
		}()

		// Handle messages from client in the main loop
		for {
			var (
				mt  int
				msg []byte
				err error
			)
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)
			// Echo the message back to client if needed
			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}
