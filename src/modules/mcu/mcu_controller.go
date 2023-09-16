package mcu

import (
	"capstonea03/be/src/contracts"
	"capstonea03/be/src/libs/db/sql"
	"capstonea03/be/src/libs/parser"
	am "capstonea03/be/src/modules/auth/auth_middleware"
	mcue "capstonea03/be/src/modules/mcu/mcu_entity"
	uc "capstonea03/be/src/modules/user/user_constant"

	"github.com/gofiber/fiber/v2"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/mcus", am.AuthGuard(uc.ROLE_ADMIN), m.getMcuList)
	m.App.Get("/api/v1/mcu/:id", am.AuthGuard(uc.ROLE_ADMIN), m.getMcu)
	m.App.Post("/api/v1/mcu", am.AuthGuard(uc.ROLE_ADMIN), m.addMcu)
	m.App.Patch("/api/v1/mcu/:id", am.AuthGuard(uc.ROLE_ADMIN), m.updateMcu)
	m.App.Delete("/api/v1/mcu/:id", am.AuthGuard(uc.ROLE_ADMIN), m.deleteMcu)
}

func (m *Module) getMcuList(c *fiber.Ctx) error {
	query := new(getMcuListReqQuery)
	if err := parser.ParseReqQuery(c, query); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	mcuListData, page, err := m.getMcuListService(&paginationOption{
		lastID: query.LastID,
		limit:  query.Limit,
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
		Data: mcuListData,
	})
}

func (m *Module) getMcu(c *fiber.Ctx) error {
	param := new(getMcuReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	mcuData, err := m.getMcuService(param.ID)
	if err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: mcuData,
	})
}

func (m *Module) addMcu(c *fiber.Ctx) error {
	req := new(addMcuReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	mcuData, err := m.addMcuService(&mcue.McuModel{
		TpsID:      req.TpsID,
		Coordinate: (*mcue.Coordinate)(req.Coordinate),
	})
	if err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(&contracts.Response{
		Data: mcuData,
	})
}

func (m *Module) updateMcu(c *fiber.Ctx) error {
	param := new(updateMcuReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	req := new(updateMcuReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	mcuData, err := m.updateMcuService(param.ID, &mcue.McuModel{
		TpsID:      req.TpsID,
		Coordinate: (*mcue.Coordinate)(req.Coordinate),
	})
	if err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: mcuData,
	})
}

func (m *Module) deleteMcu(c *fiber.Ctx) error {
	param := new(deleteMcuReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	if err := m.deleteMcuService(param.ID); err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: param.ID,
	})
}
