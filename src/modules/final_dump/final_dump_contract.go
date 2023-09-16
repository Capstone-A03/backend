package finaldump

import "github.com/google/uuid"

type getFinalDumpListReqQuery struct {
	LastID *uuid.UUID `query:"lastId"`
	Limit  *int       `query:"limit"`
}

type getFinalDumpReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type addFinalDumpReq struct {
	Name        *string        `json:"name" validate:"required,gt=0"`
	MapSectorID *uuid.UUID     `json:"mapSectorId" validate:"required"`
	Coordinate  *coordinateReq `json:"coordinate" validate:"required"`
}

type updateFinalDumpReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type updateFinalDumpReq struct {
	Name        *string        `json:"name" validate:"omitempty,gt=0"`
	MapSectorID *uuid.UUID     `json:"mapSectorId"`
	Coordinate  *coordinateReq `json:"coordinate"`
}

type deleteFinalDumpReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type coordinateReq struct {
	Latitude  *float64 `json:"latitude" validate:"required,latitude"`
	Longitude *float64 `json:"longitude" validate:"required,longitude"`
}
