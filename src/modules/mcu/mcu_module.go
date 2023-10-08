package mcu

import (
	"capstonea03/be/src/libs/db/sql"
	mcue "capstonea03/be/src/modules/mcu/mcu_entity"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	App *fiber.App
	DB  *sql.DB
}

func Load(module *Module) {
	mcue.InitRepository(module.DB)

	module.controller()
}
