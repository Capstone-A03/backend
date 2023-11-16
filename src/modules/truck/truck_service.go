package truck

import (
	"capstonea03/be/src/libs/db/sql"
	mse "capstonea03/be/src/modules/map_sector/map_sector_entity"
	te "capstonea03/be/src/modules/truck/truck_entity"

	"github.com/google/uuid"
)

type searchOption struct {
	byIsActive *bool
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

func (m *Module) getTruckListService(search *searchOption, pagination *paginationOption) (*[]*te.TruckModel, *paginationQuery, error) {
	where := new([]sql.FindAllWhere)
	limit := new(int)

	if search != nil {
		if search.byIsActive != nil {
			*where = append(*where, sql.FindAllWhere{
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
			*where = append(*where, sql.FindAllWhere{
				Where: sql.Where{
					Query: "created_at < ?",
					Args:  []interface{}{truckData.CreatedAt},
				},
				IncludeInCount: false,
			})
		}

		if pagination.limit != nil {
			*limit = *pagination.limit
		}
	}

	data, page, err := te.Repository().FindAll(&sql.FindAllOptions{
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

func (*Module) getMapSectorService(mapSectorID *uuid.UUID) (*mse.MapSectorModel, error) {
	return mse.Repository().FindOne(&sql.FindOneOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{mapSectorID},
			},
		},
	})
}

func (*Module) addTruckService(data *te.TruckModel) (*te.TruckModel, error) {
	return te.Repository().Create(data)
}

func (m *Module) updateTruckService(id *uuid.UUID, data *te.TruckModel) (*te.TruckModel, error) {
	if _, err := te.Repository().Update(data, &sql.UpdateOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	}); err != nil {
		return nil, err
	}

	return m.getTruckService(id)
}

func (*Module) deleteTruckService(id *uuid.UUID) error {
	return te.Repository().Destroy(&sql.DestroyOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}
