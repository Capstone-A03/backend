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
	Name    *string           `json:"name" validate:"required,gt=0"`
	Polygon *[]*coordinateReq `json:"polygon" validate:"required,gte=3"`
}

type updateMapSectorReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type updateMapSectorReq struct {
	Name    *string           `json:"name" validate:"omitempty,gt=0"`
	Polygon *[]*coordinateReq `json:"polygon" validate:"omitempty,gte=3"`
}

type deleteMapSectorReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type coordinateReq struct {
	Latitude  *float64 `json:"latitude" validate:"required,latitude"`
	Longitude *float64 `json:"longitude" validate:"required,longitude"`
}
