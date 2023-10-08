package truck

import (
	"capstonea03/be/src/libs/db/sql"
	te "capstonea03/be/src/modules/truck/truck_entity"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	App *fiber.App
	DB  *sql.DB
}

func Load(module *Module) {
	te.InitRepository(module.DB)

	module.controller()
}
