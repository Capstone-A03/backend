package logreport

import "capstonea03/be/src/libs/db/mongo"

type getLogReportListReqQuery struct {
	LastID *mongo.ObjectID `query:"lastId"`
	Limit  *int            `query:"limit"`
}

type getLogReportReqParam struct {
	ID *mongo.ObjectID `params:"id" validate:"required"`
}

type addLogReportReq struct {
	ReporterName *string        `json:"reporterName" validate:"required,gt=0"`
	MediaID      *string        `json:"mediaId" validate:"required,gt=0"`
	Coordinate   *coordinateReq `json:"coordinate" validate:"required"`
	Type         *string        `json:"type" validate:"required,gt=0"`
	Description  *string        `json:"description" validate:"required,gt=0"`
	Status       *string        `json:"status" validate:"required,gt=0"`
}

type updateLogReportReqParam struct {
	ID *mongo.ObjectID `params:"id" validate:"required"`
}

type updateLogReportReq struct {
	ReporterName *string        `json:"reporterName" validate:"omitempty,gt=0"`
	MediaID      *string        `json:"mediaId" validate:"omitempty,gt=0"`
	Coordinate   *coordinateReq `json:"coordinate"`
	Type         *string        `json:"type" validate:"omitempty,gt=0"`
	Description  *string        `json:"description" validate:"omitempty,gt=0"`
	Status       *string        `json:"status" validate:"omitempty,gt=0"`
}

type coordinateReq struct {
	Latitude  *float64 `json:"latitude" validate:"required,latitude"`
	Longitude *float64 `json:"longitude" validate:"required,longitude"`
}
