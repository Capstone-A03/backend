package logdump

import (
	"capstonea03/be/src/contracts"
	"capstonea03/be/src/libs/db/mongo"
	"capstonea03/be/src/libs/parser"
	am "capstonea03/be/src/modules/auth/auth_middleware"
	lde "capstonea03/be/src/modules/log_dump/log_dump_entity"
	pn "capstonea03/be/src/modules/push_notification"
	uc "capstonea03/be/src/modules/user/user_constant"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (m *Module) controller() {
	// TODO!: Add admin auth guard
	m.App.Get("/api/v1/log-dumps", m.getLogDumpList)
	m.App.Get("/api/v1/log-dump/:id", am.AuthGuard(uc.ROLE_ADMIN), m.getLogDump)
	// TODO!: Add mcu auth guard
	m.App.Post("/api/v1/log-dump", m.addLogDump)
}

func (m *Module) getLogDumpList(c *fiber.Ctx) error {
	query := new(getLogDumpListReqQuery)
	if err := parser.ParseReqQuery(c, query); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	logDumpListData, page, err := m.getLogDumpListService(
		&searchOption{
			unique: query.Unique,
			from:   query.From,
			to:     query.To,
		},
		&paginationOption{
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

	volume, err := convertVolumeToM3(req.MeasuredVolume, req.VolumeUnit)
	if err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	logDumpData, err := m.addLogDumpService(&lde.LogDumpModel{
		DumpID:         req.DumpID,
		MeasuredVolume: volume,
	})
	if err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	go func() {
		dumpData, err := m.getDumpService(req.DumpID)
		if err != nil {
			return
		}
		capacityFilled := *logDumpData.MeasuredVolume / *dumpData.Capacity * 100
		if capacityFilled >= 70 {
			pn.Send(1, "TPS sudah hampir penuh", fmt.Sprintf("%.4f%% volume TPS %s sudah terisi", capacityFilled, *dumpData.Name), pushNotification{
				Dump:    dumpData,
				LogDump: logDumpData,
			})
		}
	}()

	return c.Status(fiber.StatusCreated).JSON(&contracts.Response{
		Data: logDumpData,
	})
}
