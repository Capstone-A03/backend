package logreport

import (
	"capstonea03/be/src/libs/db/mongo"
	lre "capstonea03/be/src/modules/log_report/log_report_entity"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	App      *fiber.App
	DBClient *mongo.Client
}

func Load(module *Module) {
	lre.InitRepository(module.DBClient)

	module.controller()
}
