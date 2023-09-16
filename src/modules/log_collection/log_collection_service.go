package logcollection

import (
	"capstonea03/be/src/libs/db/mongo"
	lce "capstonea03/be/src/modules/log_collection/log_collection_entity"
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

func (m *Module) getLogCollectionListService(pagination *paginationOption) (*[]*lce.LogCollectionModel, *paginationQuery, error) {
	where := []mongo.FindAllWhere{}
	limit := 0

	if pagination.lastID != nil && len(*pagination.lastID) > 0 {
		logCollectionData, err := m.getLogCollectionService(pagination.lastID)
		if err != nil {
			return nil, nil, err
		}
		where = append(where, mongo.FindAllWhere{
			Where: mongo.Where{{
				Key: "created_at",
				Value: mongo.Where{{
					Key:   "$lt",
					Value: logCollectionData.CreatedAt,
				}},
			}},
			IncludeInCount: false,
		})
	}

	if pagination.limit != nil && *pagination.limit > 0 {
		limit = *pagination.limit
	}

	data, page, err := lce.LogCollectionRepository().FindAll(&mongo.FindAllOptions{
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

func (*Module) getLogCollectionService(id *mongo.ObjectID) (*lce.LogCollectionModel, error) {
	return lce.LogCollectionRepository().FindOne(&mongo.FindOneOptions{
		Where: &[]mongo.Where{{
			{
				Key:   "_id",
				Value: id,
			},
		}},
	})
}

func (m *Module) addLogCollectionService(data *lce.LogCollectionModel) (*lce.LogCollectionModel, error) {
	id, err := lce.LogCollectionRepository().Create(data)
	if err != nil {
		return nil, err
	}

	return m.getLogCollectionService(id)
}
