package logcollectionentity

import (
	"capstonea03/be/src/libs/db/mongo"
	"capstonea03/be/src/libs/env"
	applogger "capstonea03/be/src/libs/logger"

	"github.com/google/uuid"
)

type LogCollectionModel struct {
	mongo.Model `bson:"inline"`
	RouteID     *mongo.ObjectID `bson:"route_id,omitempty" json:"routeId"`
	TempDumpID  *uuid.UUID      `bson:"temp_dump_id,omitempty" json:"tempDumpId"`
	Volume      *float64        `bson:"volume,omitempty" json:"volume"`
	Status      *string         `bson:"status,omitempty" json:"status"`
	Note        *string         `bson:"note,omitempty" json:"note"`
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

func LogCollectionRepository() *logCollectionDB {
	if logCollectionRepo == nil {
		logger.Panic("logCollectionRepo is nil")
	}

	return logCollectionRepo
}
