package mapsector

import (
	"capstonea03/be/src/contracts"
	"capstonea03/be/src/libs/db/sql"
	"capstonea03/be/src/libs/parser"
	am "capstonea03/be/src/modules/auth/auth_middleware"
	mse "capstonea03/be/src/modules/map_sector/map_sector_entity"
	uc "capstonea03/be/src/modules/user/user_constant"

	"github.com/gofiber/fiber/v2"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/map-sectors", am.AuthGuard(uc.ROLE_ADMIN), m.getMapSectorList)
	m.App.Get("/api/v1/map-sector/:id", am.AuthGuard(uc.ROLE_ADMIN), m.getMapSector)
	m.App.Post("/api/v1/map-sector", am.AuthGuard(uc.ROLE_ADMIN), m.addMapSector)
	m.App.Patch("/api/v1/map-sector/:id", am.AuthGuard(uc.ROLE_ADMIN), m.updateMapSector)
	m.App.Delete("/api/v1/map-sector/:id", am.AuthGuard(uc.ROLE_ADMIN), m.deleteMapSector)
}

func (m *Module) getMapSectorList(c *fiber.Ctx) error {
	query := new(getMapSectorListReqQuery)
	if err := parser.ParseReqQuery(c, query); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	mapSectorData, page, err := m.getMapSectorListService(&paginationOption{
		lastID: query.LastID,
		limit:  query.Limit,
	})
	if err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Pagination: &contracts.Pagination{
			Count: page.count,
			Limit: page.limit,
			Total: page.total,
		},
		Data: mapSectorData,
	})
}

func (m *Module) getMapSector(c *fiber.Ctx) error {
	param := new(getMapSectorReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	mapSectorData, err := m.getMapSectorService(param.ID)
	if err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: mapSectorData,
	})
}

func (m *Module) addMapSector(c *fiber.Ctx) error {
	req := new(addMapSectorReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	mapSectorData, err := m.addMapSectorService(&mse.MapSectorModel{
		Name: req.Name,
		Polygon: func() *mse.Coordinates {
			println("a", len(*req.Polygon))
			if req.Polygon == nil {
				return nil
			}
			polygon := mse.Coordinates{}
			for i := range *req.Polygon {
				polygon = append(polygon, &mse.Coordinate{
					Latitude:  (*req.Polygon)[i].Latitude,
					Longitude: (*req.Polygon)[i].Longitude,
				})
			}
			return &polygon
		}(),
	})
	if err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(&contracts.Response{
		Data: mapSectorData,
	})
}

func (m *Module) updateMapSector(c *fiber.Ctx) error {
	param := new(updateMapSectorReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	req := new(updateMapSectorReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	mapSectorData, err := m.updateMapSectorService(param.ID, &mse.MapSectorModel{
		Name: req.Name,
		Polygon: func() *mse.Coordinates {
			if req.Polygon == nil {
				return nil
			}
			polygon := mse.Coordinates{}
			for i := range *req.Polygon {
				polygon = append(polygon, &mse.Coordinate{
					Latitude:  (*req.Polygon)[i].Latitude,
					Longitude: (*req.Polygon)[i].Longitude,
				})
			}
			return &polygon
		}(),
	})
	if err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: mapSectorData,
	})
}

func (m *Module) deleteMapSector(c *fiber.Ctx) error {
	param := new(deleteMapSectorReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	if err := m.deleteMapSectorService(param.ID); err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: param.ID,
	})
}
