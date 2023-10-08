package public

import (
	"capstonea03/be/src/libs/env"

	"github.com/gofiber/fiber/v2"
)

func (m *Module) controller() {
	m.App.Static("/api/v1/public", env.Get(env.APP_PUBLIC_DIRECTORY), fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        false,
		Download:      false,
		CacheDuration: -1,
	})
}
