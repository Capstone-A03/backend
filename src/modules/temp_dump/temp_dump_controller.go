package tempdump

import (
	"capstonea03/be/src/contracts"
	"capstonea03/be/src/libs/db/sql"
	"capstonea03/be/src/libs/parser"
	am "capstonea03/be/src/modules/auth/auth_middleware"
	tde "capstonea03/be/src/modules/temp_dump/temp_dump_entity"
	uc "capstonea03/be/src/modules/user/user_constant"

	"github.com/gofiber/fiber/v2"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/temp-dumps", am.AuthGuard(uc.ROLE_ADMIN), m.getTempDumpList)
	m.App.Get("/api/v1/temp-dump/:id", am.AuthGuard(uc.ROLE_ADMIN), m.getTempDump)
	m.App.Post("/api/v1/temp-dump", am.AuthGuard(uc.ROLE_ADMIN), m.addTempDump)
	m.App.Patch("/api/v1/temp-dump/:id", am.AuthGuard(uc.ROLE_ADMIN), m.updateTempDump)
	m.App.Delete("/api/v1/temp-dump/:id", am.AuthGuard(uc.ROLE_ADMIN), m.deleteTempDump)
}

func (m *Module) getTempDumpList(c *fiber.Ctx) error {
	query := new(getTempDumpListReqQuery)
	if err := parser.ParseReqQuery(c, query); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	tempDumpListData, page, err := m.getTempDumpListService(&paginationOption{
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
		Data: tempDumpListData,
	})
}

func (m *Module) getTempDump(c *fiber.Ctx) error {
	param := new(getTempDumpReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	tempDumpData, err := m.getTempDumpService(param.ID)
	if err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: tempDumpData,
	})
}

func (m *Module) addTempDump(c *fiber.Ctx) error {
	req := new(addTempDumpReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	tempDumpData, err := m.addTempDumpService(&tde.TempDumpModel{
		MapSectorID: req.MapSectorID,
		Coordinate:  (*tde.Coordinate)(req.Coordinate),
		Type:        req.Type,
		Capacity:    req.Capacity,
		Condition:   req.Condition,
	})
	if err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(&contracts.Response{
		Data: tempDumpData,
	})
}

func (m *Module) updateTempDump(c *fiber.Ctx) error {
	param := new(updateTempDumpReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	req := new(updateTempDumpReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	tempDumpData, err := m.updateTempDumpService(param.ID, &tde.TempDumpModel{
		MapSectorID: req.MapSectorID,
		Coordinate:  (*tde.Coordinate)(req.Coordinate),
		Type:        req.Type,
		Capacity:    req.Capacity,
		Condition:   req.Condition,
	})
	if err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: tempDumpData,
	})
}

func (m *Module) deleteTempDump(c *fiber.Ctx) error {
	param := new(deleteTempDumpReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	if err := m.deleteTempDumpService(param.ID); err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: param.ID,
	})
}
