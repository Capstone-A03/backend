package logdump

import (
	"capstonea03/be/src/libs/db/mongo"
	"capstonea03/be/src/libs/db/sql"
	de "capstonea03/be/src/modules/dump/dump_entity"
	lde "capstonea03/be/src/modules/log_dump/log_dump_entity"
	"capstonea03/be/src/utils"
	"sync"
	"time"

	"github.com/google/uuid"
)

type searchOption struct {
	unique *bool
	from   *time.Time
	to     *time.Time
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

func (m *Module) getLogDumpListService(search *searchOption, pagination *paginationOption) (*[]*lde.LogDumpModel, *paginationQuery, error) {
	var distinct *string = nil
	where := new([]mongo.FindAllWhere)
	limit := new(int)

	if search != nil {
		if search.unique != nil {
			distinct = utils.AsRef("_id")
		}
		if search.from != nil {
			*where = append(*where, mongo.FindAllWhere{
				Where: mongo.Where{{
					Key: "updated_at",
					Value: mongo.Where{{
						Key:   "$gte",
						Value: search.from,
					}},
				}},
			})
		}
		if search.to != nil {
			*where = append(*where, mongo.FindAllWhere{
				Where: mongo.Where{{
					Key: "updated_at",
					Value: mongo.Where{{
						Key:   "$lt",
						Value: search.to,
					}},
				}},
			})
		}
	}

	if pagination != nil {
		if !mongo.IsEmptyObjectID(pagination.lastID) {
			logDumpData, err := m.getLogDumpService(pagination.lastID)
			if err != nil {
				return nil, nil, err
			}
			*where = append(*where, mongo.FindAllWhere{
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

		if pagination.limit != nil {
			*limit = *pagination.limit
		}
	}

	if distinct != nil {
		distinctDumpID, err := lde.Repository().Distinct(&mongo.DistinctOptions{
			FieldName: "dump_id",
			Where:     where,
		})
		if err != nil {
			return nil, nil, err
		}

		distinctErrChs := make([]chan error, len(distinctDumpID))
		for i := range distinctErrChs {
			distinctErrChs[i] = make(chan error)
		}
		data := new([]*lde.LogDumpModel)
		mutex := new(sync.Mutex)
		for i, dID := range distinctDumpID {
			idx := i
			dumpID := dID
			go func() {
				_data, err := lde.Repository().FindOne(&mongo.FindOneOptions{
					Where: &[]mongo.Where{{{
						Key:   "dump_id",
						Value: dumpID,
					}}},
					Order: &[]mongo.Order{{{Key: "created_at", Value: -1}}},
				})
				if err != nil {
					distinctErrChs[idx] <- err
					return
				}
				mutex.Lock()
				*data = append(*data, _data)
				mutex.Unlock()
				distinctErrChs[idx] <- nil
			}()
		}

		for i := range distinctErrChs {
			if err := <-distinctErrChs[i]; err != nil {
				return nil, nil, err
			}
			close(distinctErrChs[i])
		}

		total := len(*data)
		return data, &paginationQuery{
			count: &total,
			limit: limit,
			total: &total,
		}, nil
	} else {
		data, page, err := lde.Repository().FindAll(&mongo.FindAllOptions{
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

func (*Module) getDumpService(id *uuid.UUID) (*de.DumpModel, error) {
	return de.Repository().FindOne(&sql.FindOneOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}

func (m *Module) addLogDumpService(data *lde.LogDumpModel) (*lde.LogDumpModel, error) {
	id, err := lde.Repository().Create(data)
	if err != nil {
		return nil, err
	}

	return m.getLogDumpService(id)
}
