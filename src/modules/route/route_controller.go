package route

import (
	"capstonea03/be/src/contracts"
	"capstonea03/be/src/libs/parser"
	am "capstonea03/be/src/modules/auth/auth_middleware"
	de "capstonea03/be/src/modules/dump/dump_entity"
	uc "capstonea03/be/src/modules/user/user_constant"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/route", am.AuthGuard(uc.ROLE_ADMIN), m.getRoute)
}

func (m *Module) getRoute(c *fiber.Ctx) error {
	query := new(getRouteReqQuery)
	if err := parser.ParseReqQuery(c, query); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	dumpListData, _, err := m.getDumpListByMapSectorIDService(query.MapSectorID)
	if err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	truckListData, _, err := m.getTruckListByMapSectorIDService(query.MapSectorID)
	if err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	finalDump := new(de.DumpModel)
	tempDumps := make([]*de.DumpModel, 0, len(*dumpListData))
	for i := range *dumpListData {
		if *(*dumpListData)[i].Type == string(de.FinalDump) {
			finalDump = (*dumpListData)[i]
		} else if *(*dumpListData)[i].Type == string(de.TempDump) {
			tempDumps = append(tempDumps, (*dumpListData)[i])
		}
	}

	routeNodes := make([]*RouteNode, 0, len(tempDumps)+1)
	routeNodes = append(routeNodes, &RouteNode{
		Coordinate: (*RouteCoordinate)(finalDump.Coordinate),
		Capacity:   finalDump.Capacity,
	})
	for i := range tempDumps {
		routeNodes = append(routeNodes, &RouteNode{
			Coordinate: (*RouteCoordinate)(tempDumps[i].Coordinate),
			Capacity:   tempDumps[i].Capacity,
		})
	}

	vehiclesCapacity := make([]float64, 0, len(*truckListData))
	for i := range *truckListData {
		vehiclesCapacity = append(vehiclesCapacity, *(*truckListData)[i].Capacity)
	}

	route := clarkeWrightSaving(&routeNodes, &vehiclesCapacity)

	routeRes := make([]*getRouteRes, 0, len(*route))
	for i := range *route {
		routeRes = append(routeRes, &getRouteRes{
			TruckID: (*truckListData)[i].ID,
		})
		dumpIDs := make([]*uuid.UUID, 0, len((*route)[i])+1)
		dumpIDs = append(dumpIDs, finalDump.ID)
		for j := range (*route)[i] {
			dumpIDs = append(dumpIDs, tempDumps[(*route)[i][j]].ID)
		}
		*routeRes[i].DumpIDs = append(*routeRes[i].DumpIDs, dumpIDs...)
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: routeRes,
	})
}
