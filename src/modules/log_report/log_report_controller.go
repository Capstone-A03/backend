package logreport

import (
	"capstonea03/be/src/contracts"
	"capstonea03/be/src/libs/db/mongo"
	"capstonea03/be/src/libs/db/sql"
	"capstonea03/be/src/libs/parser"
	lre "capstonea03/be/src/modules/log_report/log_report_entity"

	"github.com/gofiber/fiber/v2"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/log-reports", m.getLogReportList)
	m.App.Get("/api/v1/log-report/:id", m.getLogReport)
	m.App.Post("/api/v1/log-report", m.addLogReport)
	m.App.Patch("/api/v1/log-report/:id", m.updateLogReport)
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

	if _, err := m.getMediaService(req.MediaID); err != nil {
		if mongo.IsErrNoDocuments(err) {
			return contracts.NewError(fiber.ErrNotFound, "media does not exist")
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	logReportData, err := m.addLogReportService(&lre.LogReportModel{
		ReporterEmail: req.ReporterEmail,
		MediaID:       req.MediaID,
		DumpID:        req.DumpID,
		Coordinate:    (*lre.Coordinate)(req.Coordinate),
		Type:          req.Type,
		Description:   req.Description,
		Status:        req.Status,
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

	if req.MediaID != nil {
		if _, err := m.getMediaService(req.MediaID); err != nil {
			if mongo.IsErrNoDocuments(err) {
				return contracts.NewError(fiber.ErrNotFound, "media does not exist")
			}
			return contracts.NewError(fiber.ErrInternalServerError, err.Error())
		}
	}

	logReportData, err := m.updateLogReportService(param.ID, &lre.LogReportModel{
		ReporterEmail: req.ReporterEmail,
		MediaID:       req.MediaID,
		DumpID:        req.DumpID,
		Coordinate:    (*lre.Coordinate)(req.Coordinate),
		Type:          req.Type,
		Description:   req.Description,
		Status:        req.Status,
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
