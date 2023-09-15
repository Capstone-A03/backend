package auth

import (
	"time"

	"capstonea03/be/src/contracts"
	"capstonea03/be/src/libs/db/sql"
	"capstonea03/be/src/libs/hash/argon2"
	"capstonea03/be/src/libs/jwx/jwt"
	"capstonea03/be/src/libs/parser"
	ae "capstonea03/be/src/modules/auth/auth_entity"
	am "capstonea03/be/src/modules/auth/auth_middleware"
	uc "capstonea03/be/src/modules/user/user_constant"
	ue "capstonea03/be/src/modules/user/user_entity"

	"github.com/gofiber/fiber/v2"
)

func (m *Module) controller() {
	m.App.Post("/api/v1/signup", am.AuthGuard(uc.ROLE_ADMIN), m.signup)
	m.App.Post("/api/v1/signin", m.signin)
	m.App.Get("/api/v1/auth", am.AuthGuard(uc.ROLE_ADMIN, uc.ROLE_JANITOR), m.getAuth)
}

func (m *Module) signup(c *fiber.Ctx) error {
	req := new(signupReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	encodedHash, err := argon2.GetEncodedHash(req.Password)
	if err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	userData, err := m.addUserService(&ue.UserModel{
		Name:     req.Name,
		Username: req.Username,
		Password: encodedHash,
		Role:     req.Role,
	})
	if err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(&contracts.Response{
		Data: userData,
	})
}

func (m *Module) signin(c *fiber.Ctx) error {
	req := new(signinReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	userData, err := m.getUserByUsernameService(req.Username)
	if err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	isAuthorized, err := argon2.CompareStringAndEncodedHash(req.Password, userData.Password)
	if err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}
	if !isAuthorized {
		return contracts.NewError(fiber.ErrBadRequest, "incorrect username or password")
	}

	jwtToken, err := jwt.Create(&ae.JWTPayload{
		ID:   userData.ID,
		Role: userData.Role,
	})
	if err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: &userRes{
			Token: jwtToken,
			ID:    userData.ID,
			Name:  userData.Name,
			Role:  userData.Role,
		},
	})
}

func (m *Module) getAuth(c *fiber.Ctx) error {
	tokenString, err := parser.GetReqBearerToken(c)
	if err != nil {
		return contracts.NewError(fiber.ErrUnauthorized, err.Error())
	}

	token := new(ae.JWTPayload)
	if err := parser.ParseReqBearerToken(c, token); err != nil {
		return contracts.NewError(fiber.ErrUnauthorized, err.Error())
	}

	tokenExp := time.Unix(*token.Expiration, 0)
	renewToken, err := jwt.Renew[ae.JWTPayload](tokenString, &tokenExp)
	if err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}
	tokenString = renewToken

	userData, err := m.getUserService(token.ID)
	if err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: &userRes{
			Token: tokenString,
			ID:    userData.ID,
			Name:  userData.Name,
			Role:  userData.Role,
		},
	})
}
