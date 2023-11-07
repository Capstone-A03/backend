package logroute

import (
	"capstonea03/be/src/libs/db/mongo"
	lre "capstonea03/be/src/modules/log_route/log_route_entity"
	"time"

	"github.com/google/uuid"
)

type searchOption struct {
	byDriverID       *uuid.UUID
	byCreatedAtRange *[]*time.Time
}

type paginationOption struct {
	lastID *mongo.ObjectID
	limit  *int
}

type paginationQuery struct {
	count *int
	limit *int
	total *int
}

func (m *Module) getLogRouteListService(pagination *paginationOption, search *searchOption) (*[]*lre.LogRouteModel, *paginationQuery, error) {
	where := []mongo.FindAllWhere{}
	limit := 0

	if search != nil {
		if search.byDriverID != nil {
			where = append(where, mongo.FindAllWhere{
				Where: mongo.Where{{
					Key:   "driver_id",
					Value: search.byDriverID,
				}},
				IncludeInCount: true,
			})
		}

		if search.byCreatedAtRange != nil && len(*search.byCreatedAtRange) == 2 {
			where = append(where, mongo.FindAllWhere{
				Where: mongo.Where{{
					Key: "created_at",
					Value: mongo.Where{{
						Key:   "$gte",
						Value: (*search.byCreatedAtRange)[0],
					}, {
						Key:   "$lt",
						Value: (*search.byCreatedAtRange)[1],
					}},
				}},
				IncludeInCount: true,
			})
		}
	}

	if pagination != nil {
		if !mongo.IsEmptyObjectID(pagination.lastID) {
			logRouteData, err := m.getLogRouteService(pagination.lastID)
			if err != nil {
				return nil, nil, err
			}
			where = append(where, mongo.FindAllWhere{
				Where: mongo.Where{{
					Key: "created_at",
					Value: mongo.Where{{
						Key:   "$lt",
						Value: logRouteData.CreatedAt,
					}},
				}},
				IncludeInCount: false,
			})
		}

		if pagination.limit != nil && *pagination.limit > 0 {
			limit = *pagination.limit
		}
	}

	data, page, err := lre.Repository().FindAll(&mongo.FindAllOptions{
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

func (*Module) getLogRouteService(id *mongo.ObjectID) (*lre.LogRouteModel, error) {
	return lre.Repository().FindOne(&mongo.FindOneOptions{
		Where: &[]mongo.Where{{
			{
				Key:   "_id",
				Value: id,
			},
		}},
	})
}

func (m *Module) addLogRouteService(data *lre.LogRouteModel) (*lre.LogRouteModel, error) {
	id, err := lre.Repository().Create(data)
	if err != nil {
		return nil, err
	}

	return m.getLogRouteService(id)
}
