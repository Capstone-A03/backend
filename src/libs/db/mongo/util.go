package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func IsErrNoDocuments(err error) bool {
	return err == mongo.ErrNoDocuments
}

func IsEmptyObjectID(id *ObjectID) bool {
	emptyID := ObjectID{}
	return id == nil || *id == emptyID
}

func appendTimestamp[T any](data *T, setCreatedAt ...bool) error {
	docBytes, err := bson.Marshal(data)
	if err != nil {
		logger.Error(err)
		return err
	}

	docMap := map[string]interface{}{}
	if err := bson.Unmarshal(docBytes, &docMap); err != nil {
		logger.Error(err)
		return err
	}

	now := time.Now()
	if len(setCreatedAt) > 0 && setCreatedAt[0] {
		(docMap)["created_at"] = now
	}
	(docMap)["updated_at"] = now

	docBytes, err = bson.Marshal(docMap)
	if err != nil {
		logger.Error(err)
		return err
	}
	if err := bson.Unmarshal(docBytes, data); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
