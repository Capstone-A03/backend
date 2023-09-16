package logroute

import (
	"capstonea03/be/src/libs/db/mongo"

	"github.com/google/uuid"
)

type getLogRouteListReqQuery struct {
	LastID *mongo.ObjectID `query:"lastId"`
	Limit  *int            `query:"limit"`
}

type getLogRouteReqParam struct {
	ID *mongo.ObjectID `params:"id" validate:"required"`
}

type addLogRouteReq struct {
	DriverID    *uuid.UUID    `json:"driverId" validate:"required"`
	TruckID     *uuid.UUID    `json:"truckId" validate:"required"`
	TempDumpIDs *[]*uuid.UUID `json:"tempDumpIds" validate:"required,unique,gt=0"`
	FinalDumpID *uuid.UUID    `json:"finalDumpId" validate:"required"`
}
