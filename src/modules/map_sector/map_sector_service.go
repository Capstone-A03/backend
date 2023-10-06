package mapsector

import (
	"capstonea03/be/src/libs/db/sql"
	mse "capstonea03/be/src/modules/map_sector/map_sector_entity"

	"github.com/google/uuid"
)

type paginationOption struct {
	lastID *uuid.UUID
	limit  *int
}

type paginationQuery struct {
	count *int
	limit *int
	total *int
}

func (m *Module) getMapSectorListService(pagination *paginationOption) (*[]*mse.MapSectorModel, *paginationQuery, error) {
	where := []sql.FindAllWhere{}
	limit := 0

	if pagination.lastID != nil && len(*pagination.lastID) > 0 {
		mapSectorData, err := m.getMapSectorService(pagination.lastID)
		if err != nil {
			return nil, nil, err
		}
		where = append(where, sql.FindAllWhere{
			Where: sql.Where{
				Query: "created_at < ?",
				Args:  []interface{}{mapSectorData.CreatedAt},
			},
			IncludeInCount: false,
		})
	}

	if pagination.limit != nil && *pagination.limit > 0 {
		limit = *pagination.limit
	}

	data, page, err := mse.MapSectorRepository().FindAll(&sql.FindAllOptions{
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

func (*Module) getMapSectorService(id *uuid.UUID) (*mse.MapSectorModel, error) {
	return mse.MapSectorRepository().FindOne(&sql.FindOneOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}

func (*Module) addMapSectorService(data *mse.MapSectorModel) (*mse.MapSectorModel, error) {
	return mse.MapSectorRepository().Create(data)
}

func (m *Module) updateMapSectorService(id *uuid.UUID, data *mse.MapSectorModel) (*mse.MapSectorModel, error) {
	if _, err := mse.MapSectorRepository().Update(data, &sql.UpdateOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	}); err != nil {
		return nil, err
	}

	return m.getMapSectorService(id)
}

func (*Module) deleteMapSectorService(id *uuid.UUID) error {
	return mse.MapSectorRepository().Destroy(&sql.DestroyOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}
