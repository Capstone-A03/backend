package media

import (
	"capstonea03/be/src/libs/db/mongo"
	me "capstonea03/be/src/modules/media/media_entity"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	App      *fiber.App
	DBClient *mongo.Client
}

func Load(module *Module) {
	me.InitRepository(module.DBClient)

	module.controller()
}
