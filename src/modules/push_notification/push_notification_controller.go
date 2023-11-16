package pushnotification

import (
	"capstonea03/be/src/contracts"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func (m *Module) controller() {
	m.App.Use("/api/v1/push_notification", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return contracts.NewError(fiber.ErrUpgradeRequired, fiber.ErrUpgradeRequired.Error())
	})

	m.App.Get("/api/v1/push_notification", websocket.New(func(c *websocket.Conn) {
		for {
			data := <-sendCh
			bytes, err := sonic.ConfigFastest.Marshal(data)
			if err != nil {
				m.logger.Error(err)
				break
			}
			if err := c.WriteMessage(websocket.BinaryMessage, bytes); err != nil {
				m.logger.Error(err)
				break
			}
		}
	}, websocket.Config{
		RecoverHandler: func(c *websocket.Conn) {
			if err := recover(); err != nil {
				c.WriteJSON(contracts.NewError(fiber.ErrInternalServerError, fiber.ErrInternalServerError.Error()))
			}
		},
	}))

	m.App.Get("/api/v1/test_push_notification", func(c *fiber.Ctx) error {
		kindStr := c.Query("kind")
		kind, err := strconv.Atoi(kindStr)
		if err != nil {
			return contracts.NewError(fiber.ErrInternalServerError, err.Error())
		}
		title := c.Query("title")
		message := c.Query("message")

		sendCh <- getPushNotificationWsRes{
			Kind:    kind,
			Title:   title,
			Message: message,
		}

		return c.Status(fiber.StatusOK).JSON(&contracts.Response{})
	})
}
