package logdump

import (
	"capstonea03/be/src/libs/db/mongo"

	"github.com/google/uuid"
)

type getLogDumpListReqQuery struct {
	LastID *mongo.ObjectID `query:"lastId"`
	Limit  *int            `query:"limit"`
}

type getLogDumpReqParam struct {
	ID *mongo.ObjectID `params:"id" validate:"required"`
}

type addLogDumpReq struct {
	DumpID         *uuid.UUID `json:"dumpId" validate:"required"`
	MeasuredVolume *float64   `json:"measuredVolume" validate:"required"`
	MeasuredWeight *float64   `json:"measuredWeight" validate:"required"`
}
