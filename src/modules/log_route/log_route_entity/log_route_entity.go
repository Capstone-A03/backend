package logrouteentity

import (
	"capstonea03/be/src/libs/db/mongo"
	"capstonea03/be/src/libs/env"
	applogger "capstonea03/be/src/libs/logger"

	"github.com/google/uuid"
)

type LogRouteModel struct {
	mongo.Model `bson:"inline"`
	DriverID    *uuid.UUID    `bson:"driver_id,omitempty" json:"driverId,omitempty"`
	TruckID     *uuid.UUID    `bson:"truck_id,omitempty" json:"truckId,omitempty"`
	TempDumpIDs *[]*uuid.UUID `bson:"temp_dump_ids,omitempty" json:"tempDumpIds,omitempty"`
	FinalDumpID *uuid.UUID    `bson:"final_dump_id,omitempty" json:"finalDumpId,omitempty"`
}

func (LogRouteModel) DatabaseName() string {
	return env.Get(env.MONGO_DATABASE_NAME)
}

func (LogRouteModel) CollectionName() string {
	return "log_routes"
}

type logRouteDB = mongo.Service[LogRouteModel]

var logRouteRepo *logRouteDB
var logger = applogger.New("LogRouteModule")

func InitRepository(client *mongo.Client) {
	if client == nil {
		logger.Panic("client cannot be nil")
	}

	logRouteRepo = mongo.NewService[LogRouteModel](client)
}

func LogRouteRepository() *logRouteDB {
	if logRouteRepo == nil {
		logger.Panic("logRouteRepo is nil")
	}

	return logRouteRepo
}
