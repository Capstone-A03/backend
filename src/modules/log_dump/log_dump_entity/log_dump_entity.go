package logdumpentity

import (
	"capstonea03/be/src/libs/db/mongo"
	"capstonea03/be/src/libs/env"
	applogger "capstonea03/be/src/libs/logger"

	"github.com/google/uuid"
)

type LogDumpModel struct {
	mongo.Model    `bson:"inline"`
	DumpID         *uuid.UUID `bson:"dump_id,omitempty" json:"dumpId,omitempty"`
	MeasuredVolume *float64   `bson:"measured_volume,omitempty" json:"measuredVolume,omitempty"`
	MeasuredWeight *float64   `bson:"measured_weight,omitempty" json:"measuredWeight,omitempty"`
}

func (LogDumpModel) DatabaseName() string {
	return env.Get(env.MONGO_DATABASE_NAME)
}

func (LogDumpModel) CollectionName() string {
	return "log_dumps"
}

type logDumpDB = mongo.Service[LogDumpModel]

var logDumpRepo *logDumpDB
var logger = applogger.New("LogDumpModule")

func InitRepository(client *mongo.Client) {
	if client == nil {
		logger.Panic("LogDumpModule")
	}

	logDumpRepo = mongo.NewService[LogDumpModel](client)
}

func Repository() *logDumpDB {
	if logDumpRepo == nil {
		logger.Panic("logDumpRepo is nil")
	}

	return logDumpRepo
}
