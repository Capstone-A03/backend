package logcollection

import (
	"capstonea03/be/src/contracts"
	"capstonea03/be/src/libs/db/mongo"
	"capstonea03/be/src/libs/parser"
	am "capstonea03/be/src/modules/auth/auth_middleware"
	lce "capstonea03/be/src/modules/log_collection/log_collection_entity"
	uc "capstonea03/be/src/modules/user/user_constant"

	"github.com/gofiber/fiber/v2"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/log-collections", am.AuthGuard(uc.ROLE_ADMIN, uc.ROLE_JANITOR), m.getLogCollectionList)
	m.App.Get("/api/v1/log-collection/:id", am.AuthGuard(uc.ROLE_ADMIN, uc.ROLE_JANITOR), m.getLogCollection)
	m.App.Post("/api/v1/log-collection", am.AuthGuard(uc.ROLE_JANITOR), m.addLogCollection)
}

func (m *Module) getLogCollectionList(c *fiber.Ctx) error {
	query := new(getLogCollectionListReqQuery)
	if err := parser.ParseReqQuery(c, query); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	logCollectionListData, page, err := m.getLogCollectionListService(&paginationOption{
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
		Data: logCollectionListData,
	})
}

func (m *Module) getLogCollection(c *fiber.Ctx) error {
	param := new(getLogCollectionReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	logCollectionData, err := m.getLogCollectionService(param.ID)
	if err != nil {
		if mongo.IsErrNoDocuments(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: logCollectionData,
	})
}

func (m *Module) addLogCollection(c *fiber.Ctx) error {
	req := new(addLogCollectionReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	logCollectionData, err := m.addLogCollectionService(&lce.LogCollectionModel{
		RouteID: req.RouteID,
		DumpID:  req.DumpID,
		Volume:  req.Volume,
		Status:  req.Status,
		Note:    req.Note,
	})
	if err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(&contracts.Response{
		Data: logCollectionData,
	})
}
