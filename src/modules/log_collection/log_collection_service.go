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
	where := new([]mongo.FindAllWhere)
	limit := new(int)

	if pagination != nil {
		if !mongo.IsEmptyObjectID(pagination.lastID) {
			logCollectionData, err := m.getLogCollectionService(pagination.lastID)
			if err != nil {
				return nil, nil, err
			}
			*where = append(*where, mongo.FindAllWhere{
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

		if pagination.limit != nil {
			*limit = *pagination.limit
		}
	}

	data, page, err := lce.Repository().FindAll(&mongo.FindAllOptions{
		Where: where,
		Limit: limit,
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
	return lce.Repository().FindOne(&mongo.FindOneOptions{
		Where: &[]mongo.Where{{
			{
				Key:   "_id",
				Value: id,
			},
		}},
	})
}

func (m *Module) addLogCollectionService(data *lce.LogCollectionModel) (*lce.LogCollectionModel, error) {
	id, err := lce.Repository().Create(data)
	if err != nil {
		return nil, err
	}

	return m.getLogCollectionService(id)
}
