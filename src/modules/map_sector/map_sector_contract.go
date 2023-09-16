package mapsector

import "github.com/google/uuid"

type getMapSectorListReqQuery struct {
	LastID *uuid.UUID `query:"lastId"`
	Limit  *int       `query:"limit"`
}

type getMapSectorReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type addMapSectorReq struct {
	Name    *string           `json:"name" validate:"required"`
	Polygon *[]*CoordinateReq `json:"polygon" validate:"required"`
}

type updateMapSectorReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type updateMapSectorReq struct {
	Name    *string           `json:"name"`
	Polygon *[]*CoordinateReq `json:"polygon"`
}

type deleteMapSectorReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type CoordinateReq struct {
	Latitude  *float64 `json:"latitude" validate:"required"`
	Longitude *float64 `json:"longitude" validate:"required"`
}
