package logreportentity

import (
	"capstonea03/be/src/libs/db/mongo"
	"capstonea03/be/src/libs/env"
	applogger "capstonea03/be/src/libs/logger"
)

type LogReportModel struct {
	mongo.Model   `bson:"inline"`
	ReporterName  *string         `bson:"reporter_name,omitempty" json:"reporterName,omitempty"`
	ReporterEmail *string         `bson:"reporter_email,omitempty" json:"reporterEmail,omitempty"`
	MediaID       *mongo.ObjectID `bson:"media_id,omitempty" json:"mediaId,omitempty"`
	Coordinate    *Coordinate     `bson:"coordinate,omitempty" json:"coordinate,omitempty"`
	Type          *string         `bson:"type,omitempty" json:"type,omitempty"`
	Description   *string         `bson:"description,omitempty" json:"description,omitempty"`
	Status        *string         `bson:"status,omitempty" json:"status,omitempty"`
}

type Coordinate struct {
	Latitude  *float64 `bson:"latitude" json:"latitude" validate:"required,omitempty,latitude"`
	Longitude *float64 `bson:"longitude" json:"longitude" validate:"required,omitempty,longitude"`
}

func (LogReportModel) DatabaseName() string {
	return env.Get(env.MONGO_DATABASE_NAME)
}

func (LogReportModel) CollectionName() string {
	return "log_reports"
}

type logReportDB = mongo.Service[LogReportModel]

var logReportRepo *logReportDB
var logger = applogger.New("LogReportModule")

func InitRepository(client *mongo.Client) {
	if client == nil {
		logger.Panic("client cannot be nil")
	}

	logReportRepo = mongo.NewService[LogReportModel](client)
}

func Repository() *logReportDB {
	if logReportRepo == nil {
		logger.Panic("logReportRepo is nil")
	}

	return logReportRepo
}
