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
	DriverID *uuid.UUID    `json:"driverId" validate:"required"`
	TruckID  *uuid.UUID    `json:"truckId" validate:"required"`
	DumpIDs  *[]*uuid.UUID `json:"DumpIds" validate:"unique,gt=0"`
}
