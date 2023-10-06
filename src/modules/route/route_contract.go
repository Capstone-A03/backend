package route

import "github.com/google/uuid"

type getRouteReqQuery struct {
	MapSectorID *uuid.UUID `query:"mapSectorId" validate:"required"`
}

type getRouteRes struct {
	TruckID *uuid.UUID    `json:"truckId"`
	DumpIDs *[]*uuid.UUID `json:"dumpIds"`
}
