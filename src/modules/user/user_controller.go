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
	m.App.Get("/api/v1/users", m.getUserList)
	m.App.Get("/api/v1/user/:id", am.AuthGuard(uc.ROLE_ADMIN, uc.ROLE_JANITOR), m.getUser)
	m.App.Patch("/api/v1/user/:id", am.AuthGuard(uc.ROLE_ADMIN), m.updateUser)
	m.App.Delete("/api/v1/user/:id", am.AuthGuard(uc.ROLE_ADMIN), m.deleteUser)
}

func (m *Module) getUserList(c *fiber.Ctx) error {
	query := new(getUserListReqQuery)
	if err := parser.ParseReqQuery(c, query); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	userListData, page, err := m.getUserListService(&paginationOption{
		lastID: query.LastID,
		limit:  query.Limit,
	}, &searchOption{
		byRole: query.SearchByRole,
	})
	if err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Pagination: &contracts.Pagination{
			Limit: page.limit,
			Count: page.count,
			Total: page.total,
		},
		Data: userListData,
	})
}

func (m *Module) getUser(c *fiber.Ctx) error {
	token := new(ae.JWTPayload)
	if err := parser.ParseReqBearerToken(c, token); err != nil {
		return contracts.NewError(fiber.ErrUnauthorized, err.Error())
	}

	param := new(getUserReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	if *token.Role != uc.ROLE_ADMIN && *token.ID != *param.ID {
		return contracts.NewError(fiber.ErrForbidden, "you are prohibited from accessing this resource")
	}

	userData, err := m.getUserService(param.ID)
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

func (m *Module) updateUser(c *fiber.Ctx) error {
	param := new(updateUserReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

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

	userData, err := m.updateUserService(param.ID, userData)
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

	param := new(deleteUserReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	if *token.ID == *param.ID {
		if count, err := m.countUserWithAdminRoleService(); err != nil {
			return contracts.NewError(fiber.ErrInternalServerError, err.Error())
		} else if *count == 1 {
			return contracts.NewError(fiber.ErrBadRequest, "cannot delete the last admin account")
		} else if *count < 1 {
			return contracts.NewError(fiber.ErrInternalServerError, "an unexpected error occurred")
		}
	}

	if err := m.deleteUserService(param.ID); err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: param.ID,
	})
}
