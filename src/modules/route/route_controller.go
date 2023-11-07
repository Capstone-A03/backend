package route

import (
	"capstonea03/be/src/contracts"
	"capstonea03/be/src/libs/parser"
	de "capstonea03/be/src/modules/dump/dump_entity"
	te "capstonea03/be/src/modules/truck/truck_entity"
	"capstonea03/be/src/utils"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (m *Module) controller() {
	// m.App.Get("/api/v1/route", am.AuthGuard(uc.ROLE_ADMIN), m.getRoute)
	m.App.Get("/api/v1/route", m.getRoute)
}

func (m *Module) getRoute(c *fiber.Ctx) error {
	query := new(getRouteReqQuery)
	if err := parser.ParseReqQuery(c, query); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	finalDump := new(de.DumpModel)
	finalDumpErrCh := make(chan error)
	go func() {
		_finalDump, err := m.getDumpService(query.FinalDumpID)
		if err != nil {
			finalDumpErrCh <- err
			return
		}
		if *_finalDump.Type != string(de.FinalDump) {
			finalDumpErrCh <- errors.New("dump type is not final dump")
		}
		finalDump = _finalDump
		finalDumpErrCh <- nil
	}()

	tempDumpListData := new([]*de.DumpModel)
	tempDumpListDataErrCh := make(chan error)
	go func() {
		_dumpListData, _, err := m.getDumpListService(&paginationOption{
			limit: utils.AsRef(1000),
		}, &searchDumpListOption{
			mapSectorID: query.MapSectorID,
			dumpType:    utils.AsRef(string(de.TempDump)),
		})
		if err != nil {
			tempDumpListDataErrCh <- err
			return
		}
		tempDumpListData = _dumpListData
		tempDumpListDataErrCh <- nil
	}()

	truckListData := new([]*te.TruckModel)
	truckListDataErrCh := make(chan error)
	go func() {
		_truckListData, _, err := m.getTruckListService(&paginationOption{
			limit: utils.AsRef(1000),
		}, &searchTruckListOption{
			byIsActive: utils.AsRef(true),
		})
		if err != nil {
			truckListDataErrCh <- err
			return
		}
		truckListData = _truckListData
		truckListDataErrCh <- nil
	}()

	if err := <-finalDumpErrCh; err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}
	close(finalDumpErrCh)

	if err := <-tempDumpListDataErrCh; err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}
	close(tempDumpListDataErrCh)

	routeNodes := make([]*RouteNode, 0, len(*tempDumpListData)+1)
	routeNodes = append(routeNodes, &RouteNode{
		ID:                     finalDump.ID,
		Coordinate:             (*RouteCoordinate)(finalDump.Coordinate),
		RemainingWasteCapacity: finalDump.Capacity,
	})

	routeNodesErrChs := make([]chan error, len(*tempDumpListData))
	for i := range routeNodesErrChs {
		routeNodesErrChs[i] = make(chan error)
	}
	for i := range *tempDumpListData {
		idx := i
		go func() {
			logDumpData, err := m.getLogDumpByDumpIDService((*tempDumpListData)[idx].ID)
			if err != nil {
				routeNodesErrChs[idx] <- err
				return
			}
			if *logDumpData.MeasuredVolume >= *(*tempDumpListData)[idx].Capacity*70/100 {
				routeNodes = append(routeNodes, &RouteNode{
					ID:         (*tempDumpListData)[idx].ID,
					Coordinate: (*RouteCoordinate)((*tempDumpListData)[idx].Coordinate),
					RemainingWasteCapacity: func() *float64 {
						raisedCapacity := *logDumpData.MeasuredVolume * 110 / 100
						if raisedCapacity >= *(*tempDumpListData)[idx].Capacity {
							return (*tempDumpListData)[idx].Capacity
						}
						return &raisedCapacity
					}(),
				})
			}
			routeNodesErrChs[idx] <- nil
		}()
	}

	if err := <-truckListDataErrCh; err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}
	close(truckListDataErrCh)

	vehiclesCapacity := make([]*VehicleCapacity, 0, len(*truckListData))
	for i := range *truckListData {
		vehiclesCapacity = append(vehiclesCapacity, &VehicleCapacity{
			ID:       (*truckListData)[i].ID,
			Capacity: (*truckListData)[i].Capacity,
		})
	}

	for i := range routeNodesErrChs {
		if err := <-routeNodesErrChs[i]; err != nil {
			return contracts.NewError(fiber.ErrInternalServerError, err.Error())
		}
		close(routeNodesErrChs[i])
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
			dumpIDs = append(dumpIDs, (*tempDumpListData)[(*route)[i][j]].ID)
		}
		*routeRes[i].DumpIDs = append(*routeRes[i].DumpIDs, dumpIDs...)
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: routeRes,
	})
}
