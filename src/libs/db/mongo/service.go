package mongo

import (
	"capstonea03/be/src/libs/validator"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ModelI interface {
	DatabaseName() string
	CollectionName() string
}

type Service[T ModelI] struct {
	client *Client
}

type Options struct {
	UniqueFields []string
	Expiration   time.Duration
}

func NewService[T ModelI](client *Client, moptions ...*Options) *Service[T] {
	if client == nil {
		logger.Panic("db cannot be nil")
	}

	model := new(T)

	if len(moptions) > 0 {
		indexModels := func() []mongo.IndexModel {
			model := make([]mongo.IndexModel, 0)

			if len(moptions[0].UniqueFields) > 0 {
				uniqueField := bson.D{}
				for _, field := range moptions[0].UniqueFields {
					uniqueField = append(uniqueField, bson.E{Key: field, Value: 1})
				}
				model = append(model, mongo.IndexModel{
					Keys:    uniqueField,
					Options: options.Index().SetUnique(true),
				})
			}

			if moptions[0].Expiration > 0 {
				model = append(model, mongo.IndexModel{
					Keys:    bson.M{"updated_at": 1},
					Options: options.Index().SetExpireAfterSeconds(int32(moptions[0].Expiration / time.Second)),
				})
			}

			return model
		}()

		client.Database((*model).DatabaseName()).Collection((*model).CollectionName()).Indexes().CreateMany(context.TODO(), indexModels)
	}

	return &Service[T]{client: client}
}

func (s *Service[T]) Count(countOptions *CountOptions) (*int64, error) {
	model := new(T)

	coll := s.client.Database((*model).DatabaseName()).Collection((*model).CollectionName())

	where := bson.D{}
	if countOptions != nil {
		if countOptions.Where != nil {
			for i := range *countOptions.Where {
				where = append(where, (*countOptions.Where)[i]...)
			}
		}
	}

	count, err := coll.CountDocuments(context.TODO(), where)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &count, nil
}

func (s *Service[T]) FindOne(findOptions *FindOneOptions) (*T, error) {
	model := new(T)

	docStruct := new(T)

	coll := s.client.Database((*model).DatabaseName()).Collection((*model).CollectionName())
	optsFindOne := options.FindOne()

	where := bson.D{}
	if findOptions != nil {
		if findOptions.Where != nil {
			for i := range *findOptions.Where {
				where = append(where, (*findOptions.Where)[i]...)
			}
		}

		if findOptions.Order != nil {
			for i := range *findOptions.Order {
				optsFindOne = optsFindOne.SetSort((*findOptions.Order)[i])
			}
		}
	}

	if err := coll.FindOne(context.TODO(), where, optsFindOne).Decode(docStruct); err != nil {
		if !IsErrNoDocuments(err) {
			logger.Error(err)
		}
		return nil, err
	}

	return docStruct, nil
}

func (s *Service[T]) FindAll(findOptions *FindAllOptions) (*[]*T, *Pagination, error) {
	model := new(T)

	docStruct := &[]*T{}

	coll := s.client.Database((*model).DatabaseName()).Collection((*model).CollectionName())
	optsFind := options.Find()

	if findOptions != nil && findOptions.Limit != nil {
		if *findOptions.Limit > 0 && *findOptions.Limit > FindAllMaximumLimit {
			*findOptions.Limit = FindAllMaximumLimit
		} else if *findOptions.Limit == 0 {
			*findOptions.Limit = FindAllDefaultLimit
		}
	} else {
		*findOptions.Limit = FindAllDefaultLimit
	}
	if *findOptions.Limit != -1 {
		optsFind = optsFind.SetLimit(int64(*findOptions.Limit))
	}

	where := bson.D{}
	if findOptions != nil {
		if findOptions.Where != nil {
			for i := range *findOptions.Where {
				if (*findOptions.Where)[i].IncludeInCount {
					where = append(where, (*findOptions.Where)[i].Where...)
				}
			}
		}
	}

	total, err := coll.CountDocuments(context.TODO(), where)
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	if findOptions != nil && findOptions.Offset != nil && *findOptions.Offset > 0 {
		optsFind = optsFind.SetSkip(int64(*findOptions.Offset))
	} else {
		optsFind = optsFind.SetSkip(0)
	}

	if findOptions != nil {
		if findOptions.Where != nil {
			for i := range *findOptions.Where {
				if !(*findOptions.Where)[i].IncludeInCount {
					where = append(where, (*findOptions.Where)[i].Where...)
				}
			}
		}

		if findOptions.Order != nil {
			order := bson.D{}
			for i := range *findOptions.Order {
				order = append(order, (*findOptions.Order)[i]...)
			}
			optsFind = optsFind.SetSort(&order)
		}
	}

	if findOptions != nil && findOptions.Limit != nil {
		if *findOptions.Limit > FindAllMaximumLimit {
			*findOptions.Limit = FindAllMaximumLimit
		} else if *findOptions.Limit == 0 {
			*findOptions.Limit = FindAllDefaultLimit
		}
	} else {
		*findOptions.Limit = FindAllDefaultLimit
	}
	if *findOptions.Limit != -1 {
		optsFind = optsFind.SetLimit(int64(*findOptions.Limit))
	}

	cursor, err := coll.Find(context.TODO(), where, optsFind)
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	if err := cursor.All(context.TODO(), docStruct); err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	return docStruct, &Pagination{
		Limit: int(*findOptions.Limit),
		Count: len(*docStruct),
		Total: int(total),
	}, nil
}

func (s *Service[T]) Distinct(distinctOptions *DistinctOptions) ([]interface{}, error) {
	model := new(T)

	coll := s.client.Database((*model).DatabaseName()).Collection((*model).CollectionName())

	where := bson.D{}
	if distinctOptions != nil {
		if distinctOptions.Where != nil {
			for i := range *distinctOptions.Where {
				if (*distinctOptions.Where)[i].IncludeInCount {
					where = append(where, (*distinctOptions.Where)[i].Where...)
				}
			}
		}
	}

	if distinctOptions != nil {
		if distinctOptions.Where != nil {
			for i := range *distinctOptions.Where {
				if !(*distinctOptions.Where)[i].IncludeInCount {
					where = append(where, (*distinctOptions.Where)[i].Where...)
				}
			}
		}
	}

	values, err := coll.Distinct(context.TODO(), distinctOptions.FieldName, where)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return values, nil
}

func (s *Service[T]) Create(data *T) (*ObjectID, error) {
	if err := validator.Struct(data); err != nil {
		logger.Error(err)
		return nil, err
	}

	model := new(T)

	coll := s.client.Database((*model).DatabaseName()).Collection((*model).CollectionName())

	appendTimestamp(data, true)

	result, err := coll.InsertOne(context.TODO(), data)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	id := result.InsertedID.(ObjectID)

	return &id, nil
}

func (s *Service[T]) Update(data *T, updateOptions *UpdateOptions) error {
	if err := validator.Struct(data); err != nil {
		logger.Error(err)
		return err
	}

	model := new(T)

	coll := s.client.Database((*model).DatabaseName()).Collection((*model).CollectionName())

	appendTimestamp(data)

	where := bson.D{}
	if updateOptions != nil {
		if updateOptions.Where != nil {
			for i := range *updateOptions.Where {
				where = append(where, (*updateOptions.Where)[i]...)
			}
		}
	}

	result, err := coll.UpdateOne(context.TODO(), where, bson.D{{Key: "$set", Value: data}})
	if err != nil {
		logger.Error(err)
		return err
	}
	if result == nil || result.ModifiedCount == 0 {
		logger.Error("failed to update the document")
		return mongo.ErrNoDocuments
	}

	return nil
}

func (s *Service[T]) Destroy(destroyOptions *DestroyOptions) error {
	model := new(T)

	coll := s.client.Database((*model).DatabaseName()).Collection((*model).CollectionName())

	where := bson.D{}
	if destroyOptions != nil {
		if destroyOptions.Where != nil {
			for i := range *destroyOptions.Where {
				where = append(where, (*destroyOptions.Where)[i]...)
			}
		}
	}

	result, err := coll.DeleteOne(context.TODO(), where)
	if err != nil {
		logger.Error(err)
		return err
	}
	if result == nil || result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
