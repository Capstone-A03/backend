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
	TpsID     *uuid.UUID `json:"tpsId" validate:"required"`
	Latitude  *float64   `json:"latitude" validate:"required"`
	Longitude *float64   `json:"longitude" validate:"required"`
}

type updateMcuReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type updateMcuReq struct {
	TpsID     *uuid.UUID `json:"tpsId"`
	Latitude  *float64   `json:"latitude"`
	Longitude *float64   `json:"longitude"`
}

type deleteMcuReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}
