package logreport

import (
	"capstonea03/be/src/contracts"
	"capstonea03/be/src/libs/db/mongo"
	"capstonea03/be/src/libs/db/sql"
	"capstonea03/be/src/libs/parser"
	am "capstonea03/be/src/modules/auth/auth_middleware"
	lre "capstonea03/be/src/modules/log_report/log_report_entity"
	uc "capstonea03/be/src/modules/user/user_constant"

	"github.com/gofiber/fiber/v2"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/log-reports", am.AuthGuard(uc.ROLE_ADMIN), m.getLogReportList)
	m.App.Get("/api/v1/log-report/:id", am.AuthGuard(uc.ROLE_ADMIN), m.getLogReport)
	m.App.Post("/api/v1/log-report", am.AuthGuard(uc.ROLE_ADMIN), m.addLogReport)
	m.App.Patch("/api/v1/log-report/:id", am.AuthGuard(uc.ROLE_ADMIN), m.updateLogReport)
}

func (m *Module) getLogReportList(c *fiber.Ctx) error {
	query := new(getLogReportListReqQuery)
	if err := parser.ParseReqQuery(c, query); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	logReportListData, page, err := m.getLogReportListService(&paginationOption{
		lastID: query.LastID,
		limit:  query.Limit,
	})
	if err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Pagination: &contracts.Pagination{
			Limit: page.limit,
			Count: page.count,
			Total: page.total,
		},
		Data: logReportListData,
	})
}

func (m *Module) getLogReport(c *fiber.Ctx) error {
	param := new(getLogReportReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	logReportData, err := m.getLogReportService(param.ID)
	if err != nil {
		if mongo.IsErrNoDocuments(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: logReportData,
	})
}

func (m *Module) addLogReport(c *fiber.Ctx) error {
	req := new(addLogReportReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	logReportData, err := m.addLogReportService(&lre.LogReportModel{
		ReporterName: req.ReporterName,
		MediaID:      req.MediaID,
		Coordinate:   (*lre.Coordinate)(req.Coordinate),
		Type:         req.Type,
		Description:  req.Description,
		Status:       req.Status,
	})
	if err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(&contracts.Response{
		Data: logReportData,
	})
}

func (m *Module) updateLogReport(c *fiber.Ctx) error {
	param := new(updateLogReportReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	req := new(updateLogReportReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	logReportData, err := m.updateLogReportService(param.ID, &lre.LogReportModel{
		ReporterName: req.ReporterName,
		MediaID:      req.MediaID,
		Coordinate:   (*lre.Coordinate)(req.Coordinate),
		Type:         req.Type,
		Description:  req.Description,
		Status:       req.Status,
	})
	if err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: logReportData,
	})
}
