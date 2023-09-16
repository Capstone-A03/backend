package finaldump

import (
	"capstonea03/be/src/libs/db/sql"
	fde "capstonea03/be/src/modules/final_dump/final_dump_entity"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	App *fiber.App
	DB  *sql.DB
}

func Load(module *Module) {
	fde.InitRepository(module.DB)

	module.controller()
}
