package dump

import (
	"capstonea03/be/src/libs/db/sql"
	de "capstonea03/be/src/modules/dump/dump_entity"
	mse "capstonea03/be/src/modules/map_sector/map_sector_entity"

	"github.com/google/uuid"
)

type searchOption struct {
	mapSectorID *uuid.UUID
	dumpType    *string
}

type searchMapSectorOption struct {
	mapSectorID *uuid.UUID
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

func (*Module) existsDumpService(search *searchOption) (*bool, error) {
	where := []sql.Where{}

	if search != nil {
		if search.mapSectorID != nil {
			where = append(where, sql.Where{
				Query: "map_sector_id = ?",
				Args:  []interface{}{search.mapSectorID},
			})
		}

		if search.dumpType != nil {
			where = append(where, sql.Where{
				Query: "type = ?",
				Args:  []interface{}{search.dumpType},
			})
		}
	}

	return de.Repository().Exists(&sql.ExistsOptions{
		Where: &where,
	})
}

func (*Module) countMapSectorService(search *searchMapSectorOption) (*int64, error) {
	where := []sql.Where{}

	if search != nil {
		if search.mapSectorID != nil {
			where = append(where, sql.Where{
				Query: "id = ?",
				Args:  []interface{}{search.mapSectorID},
			})
		}
	}

	return mse.Repository().Count(&sql.CountOptions{
		Where: &where,
	})
}

func (m *Module) getDumpListService(search *searchOption, pagination *paginationOption) (*[]*de.DumpModel, *paginationQuery, error) {
	where := new([]sql.FindAllWhere)
	limit := new(int)

	if search != nil {
		if search.mapSectorID != nil {
			*where = append(*where, sql.FindAllWhere{
				Where: sql.Where{
					Query: "map_sector_id = ?",
					Args:  []interface{}{search.mapSectorID},
				},
				IncludeInCount: true,
			})
		}

		if search.dumpType != nil {
			*where = append(*where, sql.FindAllWhere{
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
			*where = append(*where, sql.FindAllWhere{
				Where: sql.Where{
					Query: "created_at < ?",
					Args:  []interface{}{mcuData.CreatedAt},
				},
				IncludeInCount: false,
			})
		}

		if pagination.limit != nil {
			*limit = *pagination.limit
		}
	}

	data, page, err := de.Repository().FindAll(&sql.FindAllOptions{
		Where: where,
		Limit: limit,
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

func (*Module) addDumpService(data *de.DumpModel) (*de.DumpModel, error) {
	return de.Repository().Create(data)
}

func (m *Module) updateDumpService(id *uuid.UUID, data *de.DumpModel) (*de.DumpModel, error) {
	if _, err := de.Repository().Update(data, &sql.UpdateOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	}); err != nil {
		return nil, err
	}

	return m.getDumpService(id)
}

func (*Module) deleteDumpService(id *uuid.UUID) error {
	return de.Repository().Destroy(&sql.DestroyOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}
