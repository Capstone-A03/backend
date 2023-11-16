package pushnotification

import (
	gs "capstonea03/be/src/libs/gracefulshutdown"
	"capstonea03/be/src/libs/logger"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	App    *fiber.App
	logger *logger.Logger
}

func Load(module *Module) {
	sendCh = make(chan getPushNotificationWsRes)
	gs.Add(gs.FnRunInShutdown{
		FnDescription: "closing push notification send channel",
		Fn: func() {
			close(sendCh)
		},
	})

	module.logger = logger.New("PushNotification")

	module.controller()

}
