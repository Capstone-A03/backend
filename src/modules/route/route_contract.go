package route

import (
	de "capstonea03/be/src/modules/dump/dump_entity"
	te "capstonea03/be/src/modules/truck/truck_entity"

	"github.com/google/uuid"
)

type getRouteReqQuery struct {
	FinalDumpID *uuid.UUID `query:"finalDumpId" validate:"required"`
	MapSectorID *uuid.UUID `query:"mapSectorID"`
}

type getRouteRes struct {
	Truck *te.TruckModel   `json:"truck"`
	Dumps *[]*de.DumpModel `json:"dumps"`
}
