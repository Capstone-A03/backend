package truck

import (
	"capstonea03/be/src/libs/db/sql"
	te "capstonea03/be/src/modules/truck/truck_entity"

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

func (m *Module) getTruckListService(pagination *paginationOption) (*[]*te.TruckModel, *paginationQuery, error) {
	where := []sql.FindAllWhere{}
	limit := 0

	if pagination.lastID != nil && len(*pagination.lastID) > 0 {
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

	data, page, err := te.TruckRepository().FindAll(&sql.FindAllOptions{
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
	return te.TruckRepository().FindOne(&sql.FindOneOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}

func (*Module) addTruckService(data *te.TruckModel) (*te.TruckModel, error) {
	return te.TruckRepository().Create(data)
}

func (*Module) updateTruckService(id *uuid.UUID, data *te.TruckModel) (*te.TruckModel, error) {
	where := []sql.Where{
		{
			Query: "id = ?",
			Args:  []interface{}{id},
		},
	}

	if _, err := te.TruckRepository().Update(data, &sql.UpdateOptions{
		Where: &where,
	}); err != nil {
		return nil, err
	}

	data, err := te.TruckRepository().FindOne(&sql.FindOneOptions{
		Where: &where,
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (*Module) deleteTruckService(id *uuid.UUID) error {
	return te.TruckRepository().Destroy(&te.TruckModel{
		Model: sql.Model{
			ID: id,
		},
	})
}
