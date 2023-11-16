package logroute

import (
	"capstonea03/be/src/contracts"
	"capstonea03/be/src/libs/db/mongo"
	"capstonea03/be/src/libs/parser"
	am "capstonea03/be/src/modules/auth/auth_middleware"
	lre "capstonea03/be/src/modules/log_route/log_route_entity"
	uc "capstonea03/be/src/modules/user/user_constant"

	"github.com/gofiber/fiber/v2"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/log-routes", am.AuthGuard(uc.ROLE_ADMIN, uc.ROLE_JANITOR), m.getLogRouteList)
	m.App.Get("/api/v1/log-route/:id", am.AuthGuard(uc.ROLE_ADMIN, uc.ROLE_JANITOR), m.getLogRoute)
	m.App.Post("/api/v1/log-route", am.AuthGuard(uc.ROLE_ADMIN), m.addLogRoute)
}

func (m *Module) getLogRouteList(c *fiber.Ctx) error {
	query := new(getLogRouteListReqQuery)
	if err := parser.ParseReqQuery(c, query); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	logRouteListData, page, err := m.getLogRouteListService(&searchOption{
		byDriverID:       query.DriverID,
		byCreatedAtRange: query.CreatedAtRange,
	}, &paginationOption{
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
		Data: logRouteListData,
	})
}

func (m *Module) getLogRoute(c *fiber.Ctx) error {
	param := new(getLogRouteReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	logRouteData, err := m.getLogRouteService(param.ID)
	if err != nil {
		if mongo.IsErrNoDocuments(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: logRouteData,
	})
}

func (m *Module) addLogRoute(c *fiber.Ctx) error {
	req := new(addLogRouteReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	logRouteData, err := m.addLogRouteService(&lre.LogRouteModel{
		DriverID: req.DriverID,
		TruckID:  req.TruckID,
		DumpIDs:  req.DumpIDs,
	})
	if err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(&contracts.Response{
		Data: logRouteData,
	})
}
