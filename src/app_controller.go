package main

import (
	"fmt"

	"capstonea03/be/src/contracts"
	"capstonea03/be/src/libs/env"

	"github.com/gofiber/fiber/v2"
)

func (m *module) controller() {
	m.app.Get("/", m.rootController)
}

func (*module) rootController(c *fiber.Ctx) error {
	return c.JSON(&contracts.Response{
		Data: fmt.Sprintf("%s is running", env.Get(env.APP_NAME)),
	})
}
