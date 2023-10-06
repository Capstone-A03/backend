package route

import (
	"capstonea03/be/src/libs/db/sql"
	de "capstonea03/be/src/modules/dump/dump_entity"
	te "capstonea03/be/src/modules/truck/truck_entity"
	"capstonea03/be/src/utils"

	"github.com/google/uuid"
)

func (*Module) getDumpListByMapSectorIDService(mapSectorID *uuid.UUID) (*[]*de.DumpModel, *sql.Pagination, error) {
	return de.DumpRepository().FindAll(&sql.FindAllOptions{
		Where: &[]sql.FindAllWhere{
			{
				Where: sql.Where{
					Query: "id = ?",
					Args:  []interface{}{mapSectorID},
				},
				IncludeInCount: true,
			},
		},
		Limit: utils.AsRef(sql.FindAllMaximumLimit),
	})
}

func (*Module) getTruckListByMapSectorIDService(mapSectorID *uuid.UUID) (*[]*te.TruckModel, *sql.Pagination, error) {
	return te.TruckRepository().FindAll(&sql.FindAllOptions{
		Where: &[]sql.FindAllWhere{
			{
				Where: sql.Where{
					Query: "? = ANY(map_sector_ids)",
					Args:  []interface{}{mapSectorID},
				},
			},
		},
	})
}
