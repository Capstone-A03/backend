package user

import (
	"capstonea03/be/src/libs/db/sql"
	ue "capstonea03/be/src/modules/user/user_entity"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	App *fiber.App
	DB  *sql.DB
}

func Load(module *Module) {
	ue.RegisterRoleValidation()
	ue.InitRepository(module.DB)
	ue.CreateInitialUser()

	module.controller()
}
