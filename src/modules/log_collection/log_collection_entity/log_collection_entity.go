package logcollectionentity

import (
	"capstonea03/be/src/libs/db/mongo"
	"capstonea03/be/src/libs/env"
	applogger "capstonea03/be/src/libs/logger"

	"github.com/google/uuid"
)

type LogCollectionModel struct {
	mongo.Model `bson:"inline"`
	RouteID     *mongo.ObjectID `bson:"route_id,omitempty" json:"routeId,omitempty"`
	DumpID      *uuid.UUID      `bson:"dump_id,omitempty" json:"dumpId,omitempty"`
	Volume      *float64        `bson:"volume,omitempty" json:"volume,omitempty"`
	Status      *string         `bson:"status,omitempty" json:"status,omitempty"`
	Note        *string         `bson:"note,omitempty" json:"note,omitempty"`
}

func (LogCollectionModel) DatabaseName() string {
	return env.Get(env.MONGO_DATABASE_NAME)
}

func (LogCollectionModel) CollectionName() string {
	return "log_collections"
}

type logCollectionDB = mongo.Service[LogCollectionModel]

var logCollectionRepo *logCollectionDB
var logger = applogger.New("LogCollectionModule")

func InitRepository(client *mongo.Client) {
	if client == nil {
		logger.Panic("client cannot be nil")
	}

	logCollectionRepo = mongo.NewService[LogCollectionModel](client)
}

func Repository() *logCollectionDB {
	if logCollectionRepo == nil {
		logger.Panic("logCollectionRepo is nil")
	}

	return logCollectionRepo
}
