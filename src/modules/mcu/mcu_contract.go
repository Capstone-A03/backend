package mcu

import "github.com/google/uuid"

type getMcuListReqQuery struct {
	LastID *uuid.UUID `query:"lastId"`
	Limit  *int       `query:"limit"`
}

type getMcuReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type addMcuReq struct {
	TpsID      *uuid.UUID     `json:"tpsId" validate:"required"`
	Coordinate *CoordinateReq `json:"coordinate" validate:"required"`
}

type updateMcuReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type updateMcuReq struct {
	TpsID      *uuid.UUID     `json:"tpsId"`
	Coordinate *CoordinateReq `json:"coordinate"`
}

type deleteMcuReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type CoordinateReq struct {
	Latitude  *float64 `json:"latitude" validate:"required,latitude"`
	Longitude *float64 `json:"longitude" validate:"required,longitude"`
}
