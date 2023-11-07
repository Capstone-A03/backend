package route

import (
	"capstonea03/be/src/libs/db/mongo"
	"capstonea03/be/src/libs/db/sql"
	de "capstonea03/be/src/modules/dump/dump_entity"
	lde "capstonea03/be/src/modules/log_dump/log_dump_entity"
	te "capstonea03/be/src/modules/truck/truck_entity"

	"github.com/google/uuid"
)

type searchDumpListOption struct {
	mapSectorID *uuid.UUID
	dumpType    *string
}

type searchTruckListOption struct {
	byMapSectorID *uuid.UUID
	byIsActive    *bool
}

type paginationOption struct {
	lastID *uuid.UUID
	limit  *int
}

type paginationQuery struct {
	count *int
	limit *int
	total *int
}

func (m *Module) getDumpListService(pagination *paginationOption, search *searchDumpListOption) (*[]*de.DumpModel, *paginationQuery, error) {
	where := []sql.FindAllWhere{}
	limit := 0

	if search != nil {
		if search.mapSectorID != nil {
			where = append(where, sql.FindAllWhere{
				Where: sql.Where{
					Query: "map_sector_id = ?",
					Args:  []interface{}{search.mapSectorID},
				},
				IncludeInCount: true,
			})
		}

		if search.dumpType != nil {
			where = append(where, sql.FindAllWhere{
				Where: sql.Where{
					Query: "type = ?",
					Args:  []interface{}{search.dumpType},
				},
				IncludeInCount: true,
			})
		}
	}

	if pagination != nil {
		if pagination.lastID != nil {
			mcuData, err := m.getDumpService(pagination.lastID)
			if err != nil {
				return nil, nil, err
			}
			where = append(where, sql.FindAllWhere{
				Where: sql.Where{
					Query: "created_at < ?",
					Args:  []interface{}{mcuData.CreatedAt},
				},
				IncludeInCount: false,
			})
		}

		if pagination.limit != nil && *pagination.limit > 0 {
			limit = *pagination.limit
		}
	}

	data, page, err := de.Repository().FindAll(&sql.FindAllOptions{
		Where: &where,
		Limit: &limit,
		Order: &[]string{"created_at desc"},
	})
	if err != nil {
		return nil, nil, err
	}

	return data, &paginationQuery{
		count: &page.Count,
		limit: &page.Limit,
		total: &page.Total,
	}, nil
}

func (*Module) getDumpService(id *uuid.UUID) (*de.DumpModel, error) {
	return de.Repository().FindOne(&sql.FindOneOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}

func (*Module) getLogDumpByDumpIDService(dumpID *uuid.UUID) (*lde.LogDumpModel, error) {
	return lde.Repository().FindOne(&mongo.FindOneOptions{
		Where: &[]mongo.Where{{
			{
				Key:   "dump_id",
				Value: dumpID,
			},
		}},
	})
}

func (m *Module) getTruckListService(pagination *paginationOption, search *searchTruckListOption) (*[]*te.TruckModel, *paginationQuery, error) {
	where := []sql.FindAllWhere{}
	limit := 0

	if search != nil {
		if search.byMapSectorID != nil {
			where = append(where, sql.FindAllWhere{
				Where: sql.Where{
					Query: "? = ANY(map_sector_ids)",
					Args:  []interface{}{search.byMapSectorID},
				},
				IncludeInCount: true,
			})
		}
		if search.byIsActive != nil {
			where = append(where, sql.FindAllWhere{
				Where: sql.Where{
					Query: "is_active = ?",
					Args:  []interface{}{search.byIsActive},
				},
				IncludeInCount: true,
			})
		}
	}

	if pagination != nil {
		if pagination.lastID != nil {
			truckData, err := m.getTruckService(pagination.lastID)
			if err != nil {
				return nil, nil, err
			}
			where = append(where, sql.FindAllWhere{
				Where: sql.Where{
					Query: "created_at < ?",
					Args:  []interface{}{truckData.CreatedAt},
				},
				IncludeInCount: false,
			})
		}

		if pagination.limit != nil && *pagination.limit > 0 {
			limit = *pagination.limit
		}
	}

	data, page, err := te.Repository().FindAll(&sql.FindAllOptions{
		Where: &where,
		Limit: &limit,
		Order: &[]string{"created_at desc"},
	})
	if err != nil {
		return nil, nil, err
	}

	return data, &paginationQuery{
		count: &page.Count,
		limit: &page.Limit,
		total: &page.Total,
	}, nil
}

func (*Module) getTruckService(id *uuid.UUID) (*te.TruckModel, error) {
	return te.Repository().FindOne(&sql.FindOneOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}
