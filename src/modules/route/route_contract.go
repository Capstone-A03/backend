package route

import "github.com/google/uuid"

type getRouteReqQuery struct {
	FinalDumpID *uuid.UUID `query:"finalDumpId" validate:"required"`
	MapSectorID *uuid.UUID `query:"mapSectorID"`
}

type getRouteRes struct {
	TruckID *uuid.UUID    `json:"truckId"`
	DumpIDs *[]*uuid.UUID `json:"dumpIds"`
}
