package mcu

import (
	"capstonea03/be/src/libs/db/sql"
	mcue "capstonea03/be/src/modules/mcu/mcu_entity"

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

func (m *Module) getMcuListService(pagination *paginationOption) (*[]*mcue.McuModel, *paginationQuery, error) {
	where := new([]sql.FindAllWhere)
	limit := new(int)

	if pagination != nil {
		if pagination.lastID != nil {
			mcuData, err := m.getMcuService(pagination.lastID)
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

	data, page, err := mcue.Repository().FindAll(&sql.FindAllOptions{
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

func (*Module) getMcuService(id *uuid.UUID) (*mcue.McuModel, error) {
	return mcue.Repository().FindOne(&sql.FindOneOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}

func (*Module) addMcuService(data *mcue.McuModel) (*mcue.McuModel, error) {
	return mcue.Repository().Create(data)
}

func (m *Module) updateMcuService(id *uuid.UUID, data *mcue.McuModel) (*mcue.McuModel, error) {
	if _, err := mcue.Repository().Update(data, &sql.UpdateOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	}); err != nil {
		return nil, err
	}

	return m.getMcuService(id)
}

func (*Module) deleteMcuService(id *uuid.UUID) error {
	return mcue.Repository().Destroy(&sql.DestroyOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}
