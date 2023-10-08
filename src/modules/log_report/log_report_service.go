package logreport

import (
	"capstonea03/be/src/libs/db/mongo"
	lre "capstonea03/be/src/modules/log_report/log_report_entity"
	me "capstonea03/be/src/modules/media/media_entity"
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

func (m *Module) getLogReportListService(pagination *paginationOption) (*[]*lre.LogReportModel, *paginationQuery, error) {
	where := []mongo.FindAllWhere{}
	limit := 0

	if pagination != nil {
		if !mongo.IsEmptyObjectID(pagination.lastID) {
			logReportData, err := m.getLogReportService(pagination.lastID)
			if err != nil {
				return nil, nil, err
			}
			where = append(where, mongo.FindAllWhere{
				Where: mongo.Where{{
					Key: "created_at",
					Value: mongo.Where{{
						Key:   "$lt",
						Value: logReportData.CreatedAt,
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

func (*Module) getLogReportService(id *mongo.ObjectID) (*lre.LogReportModel, error) {
	return lre.Repository().FindOne(&mongo.FindOneOptions{
		Where: &[]mongo.Where{{
			{
				Key:   "_id",
				Value: id,
			},
		}},
	})
}

func (*Module) getMediaService(mediaID *mongo.ObjectID) (*me.MediaModel, error) {
	return me.Repository().FindOne(&mongo.FindOneOptions{
		Where: &[]mongo.Where{{
			{
				Key:   "_id",
				Value: mediaID,
			},
		}},
	})
}

func (m *Module) addLogReportService(data *lre.LogReportModel) (*lre.LogReportModel, error) {
	id, err := lre.Repository().Create(data)
	if err != nil {
		return nil, err
	}

	return m.getLogReportService(id)
}

func (m *Module) updateLogReportService(id *mongo.ObjectID, data *lre.LogReportModel) (*lre.LogReportModel, error) {
	if err := lre.Repository().Update(data, &mongo.UpdateOptions{
		Where: &[]mongo.Where{{
			{
				Key:   "_id",
				Value: id,
			},
		}},
	}); err != nil {
		return nil, err
	}

	return m.getLogReportService(id)
}
