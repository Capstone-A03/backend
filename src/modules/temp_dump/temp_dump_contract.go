package tempdump

import "github.com/google/uuid"

type getTempDumpListReqQuery struct {
	LastID *uuid.UUID `query:"lastId"`
	Limit  *int       `query:"limit"`
}

type getTempDumpReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type addTempDumpReq struct {
	MapSectorID *uuid.UUID     `json:"mapSectorId" validate:"required"`
	Coordinate  *CoordinateReq `json:"coordinate" validate:"required"`
	Type        *string        `json:"type" validate:"required,gt=0"`
	Capacity    *float64       `json:"capacity" validate:"required,gt=0"`
	Condition   *string        `json:"condition" validate:"required,gt=0"`
}

type updateTempDumpReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type updateTempDumpReq struct {
	MapSectorID *uuid.UUID     `json:"mapSectorId"`
	Coordinate  *CoordinateReq `json:"coordinate"`
	Type        *string        `json:"type" validate:"omitempty,gt=0"`
	Capacity    *float64       `json:"capacity" validate:"omitempty,gt=0"`
	Condition   *string        `json:"condition" validate:"omitempty,gt=0"`
}

type deleteTempDumpReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type CoordinateReq struct {
	Latitude  *float64 `json:"latitude" validate:"required,latitude"`
	Longitude *float64 `json:"longitude" validate:"required,longitude"`
}
