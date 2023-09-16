package tempdump

import (
	"capstonea03/be/src/libs/db/sql"
	tde "capstonea03/be/src/modules/temp_dump/temp_dump_entity"

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

func (m *Module) getTempDumpListService(pagination *paginationOption) (*[]*tde.TempDumpModel, *paginationQuery, error) {
	where := []sql.FindAllWhere{}
	limit := 0

	if pagination.lastID != nil && len(*pagination.lastID) > 0 {
		mcuData, err := m.getTempDumpService(pagination.lastID)
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

	data, page, err := tde.TempDumpRepository().FindAll(&sql.FindAllOptions{
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

func (*Module) getTempDumpService(id *uuid.UUID) (*tde.TempDumpModel, error) {
	return tde.TempDumpRepository().FindOne(&sql.FindOneOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}

func (*Module) addTempDumpService(data *tde.TempDumpModel) (*tde.TempDumpModel, error) {
	return tde.TempDumpRepository().Create(data)
}

func (*Module) updateTempDumpService(id *uuid.UUID, data *tde.TempDumpModel) (*tde.TempDumpModel, error) {
	where := []sql.Where{
		{
			Query: "id = ?",
			Args:  []interface{}{id},
		},
	}

	if _, err := tde.TempDumpRepository().Update(data, &sql.UpdateOptions{
		Where: &where,
	}); err != nil {
		return nil, err
	}

	return tde.TempDumpRepository().FindOne(&sql.FindOneOptions{
		Where: &where,
	})
}

func (*Module) deleteTempDumpService(id *uuid.UUID) error {
	return tde.TempDumpRepository().Destroy(&tde.TempDumpModel{
		Model: sql.Model{
			ID: id,
		},
	})
}
