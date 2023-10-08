package mapsector

import (
	"capstonea03/be/src/libs/db/sql"
	mse "capstonea03/be/src/modules/map_sector/map_sector_entity"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	App *fiber.App
	DB  *sql.DB
}

func Load(module *Module) {
	mse.InitRepository(module.DB)

	module.controller()
}
