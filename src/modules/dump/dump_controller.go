package dump

import (
	"capstonea03/be/src/contracts"
	"capstonea03/be/src/libs/db/sql"
	"capstonea03/be/src/libs/parser"
	am "capstonea03/be/src/modules/auth/auth_middleware"
	de "capstonea03/be/src/modules/dump/dump_entity"
	uc "capstonea03/be/src/modules/user/user_constant"
	"errors"

	"github.com/gofiber/fiber/v2"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/dumps", m.getDumpList)
	m.App.Get("/api/v1/dump/:id", m.getDump)
	m.App.Post("/api/v1/dump", am.AuthGuard(uc.ROLE_ADMIN), m.addDump)
	m.App.Patch("/api/v1/dump/:id", am.AuthGuard(uc.ROLE_ADMIN), m.updateDump)
	m.App.Delete("/api/v1/dump/:id", am.AuthGuard(uc.ROLE_ADMIN), m.deleteDump)
}

func (m *Module) getDumpList(c *fiber.Ctx) error {
	query := new(getDumpListReqQuery)
	if err := parser.ParseReqQuery(c, query); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	dumpListData, page, err := m.getDumpListService(&searchOption{
		mapSectorID: query.SearchByMapSectorID,
		dumpType:    query.SearchByType,
	}, &paginationOption{
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
		Data: dumpListData,
	})
}

func (m *Module) getDump(c *fiber.Ctx) error {
	param := new(getDumpReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	dumpData, err := m.getDumpService(param.ID)
	if err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: dumpData,
	})
}

func (m *Module) addDump(c *fiber.Ctx) error {
	req := new(addDumpReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	if req.MapSectorID != nil {
		countMapSector, err := m.countMapSectorService(&searchMapSectorOption{
			mapSectorID: req.MapSectorID,
		})
		if err != nil {
			return contracts.NewError(fiber.ErrInternalServerError, err.Error())
		}
		if *countMapSector == 0 {
			return contracts.NewError(fiber.ErrBadRequest, "the map sector doesn't exist")
		}
	}

	if err := func() error {
		switch *req.Type {
		case string(de.TempDump):
			return nil
		case string(de.FinalDump):
			return nil
		default:
			return errors.New("invalid dump type")
		}
	}(); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	if *req.Type == string(de.FinalDump) && req.MapSectorID != nil {
		existsFinalDump, err := m.existsDumpService(&searchOption{
			mapSectorID: req.MapSectorID,
			dumpType:    req.Type,
		})
		if err != nil {
			return contracts.NewError(fiber.ErrInternalServerError, err.Error())
		}
		if *existsFinalDump {
			return contracts.NewError(fiber.ErrBadRequest, "the final dump is already exists")
		}
	}

	dumpData, err := m.addDumpService(&de.DumpModel{
		Name:        req.Name,
		MapSectorID: req.MapSectorID,
		Coordinate:  (*de.Coordinate)(req.Coordinate),
		Type:        req.Type,
		Capacity:    req.Capacity,
		Condition:   req.Condition,
	})
	if err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(&contracts.Response{
		Data: dumpData,
	})
}

func (m *Module) updateDump(c *fiber.Ctx) error {
	param := new(updateDumpReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	req := new(updateDumpReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	if req.MapSectorID != nil {
		countMapSector, err := m.countMapSectorService(&searchMapSectorOption{
			mapSectorID: req.MapSectorID,
		})
		if err != nil {
			return contracts.NewError(fiber.ErrInternalServerError, err.Error())
		}
		if *countMapSector == 0 {
			return contracts.NewError(fiber.ErrBadRequest, "the map sector doesn't exist")
		}
	}

	if err := func() error {
		if req.Type == nil {
			return nil
		}

		switch *req.Type {
		case string(de.TempDump):
			return nil
		case string(de.FinalDump):
			return nil
		default:
			return errors.New("invalid dump type")
		}
	}(); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	dumpData, err := m.getDumpService(param.ID)
	if err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	if req.Type != nil && *req.Type == string(de.FinalDump) && (req.MapSectorID != nil || dumpData.MapSectorID != nil) {
		mapSectorID := dumpData.MapSectorID
		if req.MapSectorID != nil {
			mapSectorID = req.MapSectorID
		}
		existsFinalDump, err := m.existsDumpService(&searchOption{
			mapSectorID: mapSectorID,
			dumpType:    req.Type,
		})
		if err != nil {
			return contracts.NewError(fiber.ErrInternalServerError, err.Error())
		}
		if *existsFinalDump {
			return contracts.NewError(fiber.ErrBadRequest, "the final dump is already exists")
		}
	}

	dumpData, err = m.updateDumpService(param.ID, &de.DumpModel{
		Name:        req.Name,
		MapSectorID: req.MapSectorID,
		Coordinate:  (*de.Coordinate)(req.Coordinate),
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
		Data: dumpData,
	})
}

func (m *Module) deleteDump(c *fiber.Ctx) error {
	param := new(deleteDumpReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	if err := m.deleteDumpService(param.ID); err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: param.ID,
	})
}
