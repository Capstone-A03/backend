package dump

import (
	"capstonea03/be/src/libs/db/sql"
	de "capstonea03/be/src/modules/dump/dump_entity"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	App *fiber.App
	DB  *sql.DB
}

func Load(module *Module) {
	de.InitRepository(module.DB)

	module.controller()
}
