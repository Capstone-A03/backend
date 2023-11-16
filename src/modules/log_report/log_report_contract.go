package logreport

import (
	"capstonea03/be/src/libs/db/mongo"

	"github.com/google/uuid"
)

type getLogReportListReqQuery struct {
	LastID *mongo.ObjectID `query:"lastId"`
	Limit  *int            `query:"limit"`
}

type getLogReportReqParam struct {
	ID *mongo.ObjectID `params:"id" validate:"required"`
}

type addLogReportReq struct {
	ReporterEmail *string         `json:"reporterEmail" validate:"omitempty"`
	MediaID       *mongo.ObjectID `json:"mediaId" validate:"required"`
	DumpID        *uuid.UUID      `json:"dumpId" validate:"omitempty"`
	Coordinate    *coordinate     `json:"coordinate" validate:"required"`
	Type          *string         `json:"type" validate:"gt=0"`
	Description   *string         `json:"description" validate:"gt=0"`
	Status        *string         `json:"status" validate:"gt=0"`
}

type updateLogReportReqParam struct {
	ID *mongo.ObjectID `params:"id" validate:"required"`
}

type updateLogReportReq struct {
	ReporterEmail *string         `json:"reporterEmail" validate:"omitempty"`
	MediaID       *mongo.ObjectID `json:"mediaId" validate:"omitempty"`
	DumpID        *uuid.UUID      `json:"dumpId" validate:"omitempty"`
	Coordinate    *coordinate     `json:"coordinate" validate:"omitempty"`
	Type          *string         `json:"type" validate:"omitempty,gt=0"`
	Description   *string         `json:"description" validate:"omitempty,gt=0"`
	Status        *string         `json:"status" validate:"omitempty,gt=0"`
}

type coordinate struct {
	Latitude  *float64 `json:"latitude" validate:"required,latitude"`
	Longitude *float64 `json:"longitude" validate:"required,longitude"`
}
