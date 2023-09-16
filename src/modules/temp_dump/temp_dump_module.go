package tempdump

import (
	"capstonea03/be/src/libs/db/sql"
	tde "capstonea03/be/src/modules/temp_dump/temp_dump_entity"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	App *fiber.App
	DB  *sql.DB
}

func Load(module *Module) {
	tde.InitRepository(module.DB)

	module.controller()
}
