package logroute

import (
	"capstonea03/be/src/libs/db/mongo"
	"time"

	"github.com/google/uuid"
)

type getLogRouteListReqQuery struct {
	DriverID       *uuid.UUID      `query:"driverId"`
	CreatedAtRange *[]*time.Time   `query:"createdAtRange" validate:"omitempty,len=2"`
	LastID         *mongo.ObjectID `query:"lastId"`
	Limit          *int            `query:"limit"`
}

type getLogRouteReqParam struct {
	ID *mongo.ObjectID `params:"id" validate:"required"`
}

type addLogRouteReq struct {
	DriverID *uuid.UUID    `json:"driverId" validate:"required"`
	TruckID  *uuid.UUID    `json:"truckId" validate:"required"`
	DumpIDs  *[]*uuid.UUID `json:"dumpIds" validate:"unique,gt=0"`
}
