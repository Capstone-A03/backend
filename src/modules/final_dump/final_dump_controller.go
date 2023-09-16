package finaldump

import (
	"capstonea03/be/src/contracts"
	"capstonea03/be/src/libs/db/sql"
	"capstonea03/be/src/libs/parser"
	am "capstonea03/be/src/modules/auth/auth_middleware"
	fde "capstonea03/be/src/modules/final_dump/final_dump_entity"
	uc "capstonea03/be/src/modules/user/user_constant"

	"github.com/gofiber/fiber/v2"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/final-dumps", am.AuthGuard(uc.ROLE_ADMIN), m.getFinalDumpList)
	m.App.Get("/api/v1/final-dump/:id", am.AuthGuard(uc.ROLE_ADMIN), m.getFinalDump)
	m.App.Post("/api/v1/final-dump", am.AuthGuard(uc.ROLE_ADMIN), m.addFinalDump)
	m.App.Patch("/api/v1/final-dump/:id", am.AuthGuard(uc.ROLE_ADMIN), m.updateFinalDump)
	m.App.Delete("/api/v1/final-dump/:id", am.AuthGuard(uc.ROLE_ADMIN), m.deleteFinalDump)
}

func (m *Module) getFinalDumpList(c *fiber.Ctx) error {
	query := new(getFinalDumpListReqQuery)
	if err := parser.ParseReqQuery(c, query); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	finalDumpListData, page, err := m.getFinalDumpListService(&paginationOption{
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
		Data: finalDumpListData,
	})
}

func (m *Module) getFinalDump(c *fiber.Ctx) error {
	param := new(getFinalDumpReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	finalDumpData, err := m.getFinalDumpService(param.ID)
	if err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: finalDumpData,
	})
}

func (m *Module) addFinalDump(c *fiber.Ctx) error {
	req := new(addFinalDumpReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	finalDumpData, err := m.addFinalDumpService(&fde.FinalDumpModel{
		Name:        req.Name,
		MapSectorID: req.MapSectorID,
		Coordinate:  (*fde.Coordinate)(req.Coordinate),
	})
	if err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(&contracts.Response{
		Data: finalDumpData,
	})
}

func (m *Module) updateFinalDump(c *fiber.Ctx) error {
	param := new(updateFinalDumpReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	req := new(updateFinalDumpReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	finalDumpData, err := m.updateFinalDumpService(param.ID, &fde.FinalDumpModel{
		Name:        req.Name,
		MapSectorID: req.MapSectorID,
		Coordinate:  (*fde.Coordinate)(req.Coordinate),
	})
	if err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: finalDumpData,
	})
}

func (m *Module) deleteFinalDump(c *fiber.Ctx) error {
	param := new(deleteFinalDumpReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	if err := m.deleteFinalDumpService(param.ID); err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: param.ID,
	})
}
