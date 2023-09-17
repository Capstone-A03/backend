package logdump

import (
	"capstonea03/be/src/contracts"
	"capstonea03/be/src/libs/db/mongo"
	"capstonea03/be/src/libs/parser"
	am "capstonea03/be/src/modules/auth/auth_middleware"
	lde "capstonea03/be/src/modules/log_dump/log_dump_entity"
	uc "capstonea03/be/src/modules/user/user_constant"

	"github.com/gofiber/fiber/v2"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/log-dumps", am.AuthGuard(uc.ROLE_ADMIN), m.getLogDumpList)
	m.App.Get("/api/v1/log-dump/:id", am.AuthGuard(uc.ROLE_ADMIN), m.getLogDump)
	m.App.Post("/api/v1/log-dump", am.AuthGuard(uc.ROLE_ADMIN), m.addLogDump)
}

func (m *Module) getLogDumpList(c *fiber.Ctx) error {
	query := new(getLogDumpListReqQuery)
	if err := parser.ParseReqQuery(c, query); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	logDumpListData, page, err := m.getLogDumpListService(&paginationOption{
		lastID: query.LastID,
		limit:  query.Limit,
	})
	if err != nil {
		if mongo.IsErrNoDocuments(err) {
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
		Data: logDumpListData,
	})
}

func (m *Module) getLogDump(c *fiber.Ctx) error {
	param := new(getLogDumpReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	logDumpData, err := m.getLogDumpService(param.ID)
	if err != nil {
		if mongo.IsErrNoDocuments(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: logDumpData,
	})
}

func (m *Module) addLogDump(c *fiber.Ctx) error {
	req := new(addLogDumpReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	logDumpData, err := m.addLogDumpService(&lde.LogDumpModel{
		DumpID:         req.DumpID,
		MeasuredVolume: req.MeasuredVolume,
		MeasuredWeight: req.MeasuredWeight,
	})
	if err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(&contracts.Response{
		Data: logDumpData,
	})
}
