package pushnotification

import (
	"capstonea03/be/src/contracts"
	"strconv"
	"time"

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
		query := new(getPushNotificationWsReqQuery)
		query.token = c.Query("token")
		if len(query.token) == 0 {
			errRes := contracts.NewError(fiber.ErrBadRequest, "token required")
			if err := contracts.NewWSRes(c, errRes); err != nil {
				m.logger.Error(err)
			}
			if err := c.Close(); err != nil {
				m.logger.Error(err)
			}
			return
		}

		ticker := time.NewTicker(5 * time.Second)
		tickerChan := ticker.C

		closeCh := make(chan bool)

		clientPing := false
		go func() {
			for {
				select {
				case isClose := <-closeCh:
					if isClose {
						return
					}
				default:
					if _, msg, err := c.ReadMessage(); err == nil {
						if string(msg) == "PING" {
							clientPing = true
							ticker.Reset(5 * time.Second)
						}
					} else {
						return
					}
				}

			}
		}()

		go func() {
			for {
				select {
				case data := <-sendCh:
					if err := contracts.NewWSRes(c, data); err != nil {
						m.logger.Error(err)
						closeCh <- true
					}
				case <-tickerChan:
					if !clientPing {
						closeCh <- true
						ticker.Stop()
					}
				case isClose := <-closeCh:
					if isClose {
						c.Close()
						return
					}
				}
			}
		}()

		<-closeCh
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
