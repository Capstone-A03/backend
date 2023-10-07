package truck

import (
	"capstonea03/be/src/contracts"
	"capstonea03/be/src/libs/db/sql"
	"capstonea03/be/src/libs/parser"
	am "capstonea03/be/src/modules/auth/auth_middleware"
	te "capstonea03/be/src/modules/truck/truck_entity"
	uc "capstonea03/be/src/modules/user/user_constant"

	"github.com/gofiber/fiber/v2"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/trucks", am.AuthGuard(uc.ROLE_ADMIN), m.getTruckList)
	m.App.Get("/api/v1/truck/:id", am.AuthGuard(uc.ROLE_ADMIN), m.getTruck)
	m.App.Post("/api/v1/truck", am.AuthGuard(uc.ROLE_ADMIN), m.addTruck)
	m.App.Patch("/api/v1/truck/:id", am.AuthGuard(uc.ROLE_ADMIN), m.updateTruck)
	m.App.Delete("/api/v1/truck/:id", am.AuthGuard(uc.ROLE_ADMIN), m.deleteTruck)
}

func (m *Module) getTruckList(c *fiber.Ctx) error {
	query := new(getTruckListReqQuery)
	if err := parser.ParseReqQuery(c, query); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	truckListData, page, err := m.getTruckListService(&paginationOption{
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
		Data: truckListData,
	})
}

func (m *Module) getTruck(c *fiber.Ctx) error {
	param := new(getTruckReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	truckData, err := m.getTruckService(param.ID)
	if err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: truckData,
	})
}

func (m *Module) addTruck(c *fiber.Ctx) error {
	req := new(addTruckReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	ch := make([]chan error, len(*req.MapSectorIDs))
	for i := range *req.MapSectorIDs {
		idx := i
		go func() {
			defer close(ch[idx])

			if _, err := m.getMapSectorService((*req.MapSectorIDs)[idx]); err != nil {
				ch[idx] <- err
			}
			ch[idx] <- nil
		}()
	}
	for i := range ch {
		err := <-ch[i]
		if err != nil {
			if sql.IsErrRecordNotFound(err) {
				return contracts.NewError(fiber.ErrBadRequest, err.Error())
			}
			return contracts.NewError(fiber.ErrInternalServerError, err.Error())
		}
	}

	truckData, err := m.addTruckService(&te.TruckModel{
		MapSectorIDs:    req.MapSectorIDs,
		LicensePlate:    req.LicensePlate,
		Type:            req.Type,
		Capacity:        req.Capacity,
		FuelConsumption: req.FuelConsumption,
	})
	if err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(&contracts.Response{
		Data: truckData,
	})
}

func (m *Module) updateTruck(c *fiber.Ctx) error {
	param := new(updateTruckReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	req := new(updateTruckReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	if req.MapSectorIDs != nil {
		ch := make([]chan error, len(*req.MapSectorIDs))
		for i := range *req.MapSectorIDs {
			idx := i
			go func() {
				defer close(ch[idx])

				if _, err := m.getMapSectorService((*req.MapSectorIDs)[idx]); err != nil {
					ch[idx] <- err
				}
				ch[idx] <- nil
			}()
		}
		for i := range ch {
			err := <-ch[i]
			if err != nil {
				if sql.IsErrRecordNotFound(err) {
					return contracts.NewError(fiber.ErrBadRequest, err.Error())
				}
				return contracts.NewError(fiber.ErrInternalServerError, err.Error())
			}
		}
	}

	truckData, err := m.updateTruckService(param.ID, &te.TruckModel{
		MapSectorIDs:    req.MapSectorIDs,
		LicensePlate:    req.LicensePlate,
		Type:            req.Type,
		Capacity:        req.Capacity,
		FuelConsumption: req.FuelConsumption,
	})
	if err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: truckData,
	})
}

func (m *Module) deleteTruck(c *fiber.Ctx) error {
	param := new(deleteTruckReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	if err := m.deleteTruckService(param.ID); err != nil {
		if sql.IsErrRecordNotFound(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: param.ID,
	})
}
