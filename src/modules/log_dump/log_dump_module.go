package logdump

import (
	"capstonea03/be/src/libs/db/mongo"
	lde "capstonea03/be/src/modules/log_dump/log_dump_entity"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	App      *fiber.App
	DBClient *mongo.Client
}

func Load(module *Module) {
	lde.InitRepository(module.DBClient)

	module.controller()
}
