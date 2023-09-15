package authmiddleware

import (
	"capstonea03/be/src/contracts"
	"capstonea03/be/src/libs/parser"
	ae "capstonea03/be/src/modules/auth/auth_entity"
	uc "capstonea03/be/src/modules/user/user_constant"

	"github.com/gofiber/fiber/v2"
)

func AuthGuard(role ...uc.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(role) == 0 {
			return c.Next()
		}

		token := new(ae.JWTPayload)
		if err := parser.ParseReqBearerToken(c, token); err != nil {
			return contracts.NewError(fiber.ErrUnauthorized, err.Error())
		}

		isAuthorized := false
		for i := range role {
			if role[i] == *token.Role {
				isAuthorized = true
			}
		}

		if !isAuthorized {
			return contracts.NewError(fiber.ErrForbidden, "you are prohibited from accessing this resource")
		}

		return c.Next()
	}
}
