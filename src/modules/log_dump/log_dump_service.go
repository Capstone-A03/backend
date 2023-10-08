package logdump

import (
	"capstonea03/be/src/libs/db/mongo"
	lde "capstonea03/be/src/modules/log_dump/log_dump_entity"
)

type paginationOption struct {
	lastID *mongo.ObjectID
	limit  *int
}

type paginationQuery struct {
	count *int
	limit *int
	total *int
}

func (m *Module) getLogDumpListService(pagination *paginationOption) (*[]*lde.LogDumpModel, *paginationQuery, error) {
	where := []mongo.FindAllWhere{}
	limit := 0

	if pagination != nil {
		if !mongo.IsEmptyObjectID(pagination.lastID) {
			logDumpData, err := m.getLogDumpService(pagination.lastID)
			if err != nil {
				return nil, nil, err
			}
			where = append(where, mongo.FindAllWhere{
				Where: mongo.Where{{
					Key: "created_at",
					Value: mongo.Where{{
						Key:   "$lt",
						Value: logDumpData.CreatedAt,
					}},
				}},
				IncludeInCount: false,
			})
		}

		if pagination.limit != nil && *pagination.limit > 0 {
			limit = *pagination.limit
		}
	}

	data, page, err := lde.Repository().FindAll(&mongo.FindAllOptions{
		Where: &where,
		Limit: &limit,
		Order: &[]mongo.Order{{{Key: "created_at", Value: -1}}},
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

func (*Module) getLogDumpService(id *mongo.ObjectID) (*lde.LogDumpModel, error) {
	return lde.Repository().FindOne(&mongo.FindOneOptions{
		Where: &[]mongo.Where{{
			{
				Key:   "_id",
				Value: id,
			},
		}},
	})
}

func (m *Module) addLogDumpService(data *lde.LogDumpModel) (*lde.LogDumpModel, error) {
	id, err := lde.Repository().Create(data)
	if err != nil {
		return nil, err
	}

	return m.getLogDumpService(id)
}
