package user

import (
	"capstonea03/be/src/contracts"
	"capstonea03/be/src/libs/db/sql"
	"capstonea03/be/src/libs/hash/argon2"
	"capstonea03/be/src/libs/parser"
	ae "capstonea03/be/src/modules/auth/auth_entity"
	am "capstonea03/be/src/modules/auth/auth_middleware"
	uc "capstonea03/be/src/modules/user/user_constant"
	ue "capstonea03/be/src/modules/user/user_entity"

	"github.com/gofiber/fiber/v2"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/user", am.AuthGuard(uc.ROLE_ADMIN, uc.ROLE_JANITOR), m.getUser)
	m.App.Patch("/api/v1/user", am.AuthGuard(uc.ROLE_ADMIN), m.updateUser)
	m.App.Delete("/api/v1/user", am.AuthGuard(uc.ROLE_ADMIN), m.deleteUser)
}

func (m *Module) getUser(c *fiber.Ctx) error {
	token := new(ae.JWTPayload)
	if err := parser.ParseReqBearerToken(c, token); err != nil {
		return contracts.NewError(fiber.ErrUnauthorized, err.Error())
	}

	query := new(getUserReqQuery)
	if err := parser.ParseReqQuery(c, query); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	if *token.Role != uc.ROLE_ADMIN && *token.ID != *query.ID {
		return contracts.NewError(fiber.ErrForbidden, "you are prohibited from accessing this resource")
	}

	userData, err := m.getUserService(query.ID)
	if err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrUnauthorized, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: userData,
	})
}

func (m *Module) updateUser(c *fiber.Ctx) error {
	req := new(updateUserReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	userData := &ue.UserModel{
		Name:     req.Name,
		Username: req.Username,
	}

	if req.Password != nil && len(*req.Password) > 0 {
		encodedHash, err := argon2.GetEncodedHash(req.Password)
		if err != nil {
			return contracts.NewError(fiber.ErrInternalServerError, err.Error())
		}
		userData.Password = encodedHash
	}

	userData, err := m.updateUserService(req.ID, userData)
	if err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: userData,
	})
}

func (m *Module) deleteUser(c *fiber.Ctx) error {
	token := new(ae.JWTPayload)
	if err := parser.ParseReqBearerToken(c, token); err != nil {
		return contracts.NewError(fiber.ErrUnauthorized, err.Error())
	}

	req := new(deleteUserReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	if *token.ID == *req.ID {
		if count, err := m.countUserWithAdminRoleService(); err != nil {
			return contracts.NewError(fiber.ErrInternalServerError, err.Error())
		} else if *count == 1 {
			return contracts.NewError(fiber.ErrBadRequest, "cannot delete the last admin account")
		} else if *count < 1 {
			return contracts.NewError(fiber.ErrInternalServerError, "an unexpected error occurred")
		}
	}

	if err := m.deleteUserService(req.ID); err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: req.ID,
	})
}
