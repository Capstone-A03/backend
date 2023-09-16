package logcollection

import (
	"capstonea03/be/src/libs/db/mongo"

	"github.com/google/uuid"
)

type getLogCollectionListReqQuery struct {
	LastID *mongo.ObjectID `query:"lastId"`
	Limit  *int            `query:"limit"`
}

type getLogCollectionReqParam struct {
	ID *mongo.ObjectID `params:"id" validate:"required"`
}

type addLogCollectionReq struct {
	RouteID    *mongo.ObjectID `json:"routeId" validate:"required"`
	TempDumpID *uuid.UUID      `json:"tempDumpId" validate:"required"`
	Volume     *float64        `json:"volume" validate:"required,gt=0"`
	Status     *string         `json:"status" validate:"required,gt=0"`
	Note       *string         `json:"note"`
}
