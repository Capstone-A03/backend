package logcollection

import (
	"capstonea03/be/src/libs/db/mongo"
	lce "capstonea03/be/src/modules/log_collection/log_collection_entity"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	App      *fiber.App
	DBClient *mongo.Client
}

func Load(module *Module) {
	lce.InitRepository(module.DBClient)

	module.controller()
}
