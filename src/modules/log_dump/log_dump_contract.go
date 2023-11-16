package logdump

import (
	"capstonea03/be/src/libs/db/mongo"
	"time"

	"github.com/google/uuid"
)

type getLogDumpListReqQuery struct {
	Unique *bool           `query:"unique"`
	From   *time.Time      `query:"from"`
	To     *time.Time      `query:"to"`
	LastID *mongo.ObjectID `query:"lastId"`
	Limit  *int            `query:"limit"`
}

type getLogDumpReqParam struct {
	ID *mongo.ObjectID `params:"id" validate:"required"`
}

type addLogDumpReq struct {
	DumpID         *uuid.UUID `json:"dumpId" validate:"required"`
	MeasuredVolume *float64   `json:"measuredVolume" validate:"required"`
	VolumeUnit     *string    `json:"volumeUnit"`
}
