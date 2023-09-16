package finaldump

import (
	"capstonea03/be/src/libs/db/sql"
	fde "capstonea03/be/src/modules/final_dump/final_dump_entity"

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

func (m *Module) getFinalDumpListService(pagination *paginationOption) (*[]*fde.FinalDumpModel, *paginationQuery, error) {
	where := []sql.FindAllWhere{}
	limit := 0

	if pagination.lastID != nil && len(*pagination.lastID) > 0 {
		mcuData, err := m.getFinalDumpService(pagination.lastID)
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

	data, page, err := fde.FinalDumpRepository().FindAll(&sql.FindAllOptions{
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

func (*Module) getFinalDumpService(id *uuid.UUID) (*fde.FinalDumpModel, error) {
	return fde.FinalDumpRepository().FindOne(&sql.FindOneOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}

func (*Module) addFinalDumpService(data *fde.FinalDumpModel) (*fde.FinalDumpModel, error) {
	return fde.FinalDumpRepository().Create(data)
}

func (m *Module) updateFinalDumpService(id *uuid.UUID, data *fde.FinalDumpModel) (*fde.FinalDumpModel, error) {
	if _, err := fde.FinalDumpRepository().Update(data, &sql.UpdateOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	}); err != nil {
		return nil, err
	}

	return m.getFinalDumpService(id)
}

func (*Module) deleteFinalDumpService(id *uuid.UUID) error {
	return fde.FinalDumpRepository().Destroy(&fde.FinalDumpModel{
		Model: sql.Model{
			ID: id,
		},
	})
}
