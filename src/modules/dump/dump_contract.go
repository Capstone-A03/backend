package dump

import "github.com/google/uuid"

type getDumpListReqQuery struct {
	SearchByMapSectorID *uuid.UUID `query:"mapSectorId"`
	SearchByType        *string    `query:"type"`
	LastID              *uuid.UUID `query:"lastId"`
	Limit               *int       `query:"limit"`
}

type getDumpReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type addDumpReq struct {
	Name        *string     `json:"name" validate:"gt=0"`
	MapSectorID *uuid.UUID  `json:"mapSectorId" validate:"omitempty"`
	Coordinate  *coordinate `json:"coordinate" validate:"required"`
	Type        *string     `json:"type" validate:"gt=0"`
	Capacity    *float64    `json:"capacity" validate:"gt=0"`
	Condition   *string     `json:"condition" validate:"omitempty,gt=0"`
}

type updateDumpReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type updateDumpReq struct {
	Name        *string     `json:"name" validate:"omitempty,gt=0"`
	MapSectorID *uuid.UUID  `json:"mapSectorId" validate:"omitempty"`
	Coordinate  *coordinate `json:"coordinate"`
	Type        *string     `json:"type" validate:"omitempty,gt=0"`
	Capacity    *float64    `json:"capacity" validate:"omitempty,gt=0"`
	Condition   *string     `json:"condition" validate:"omitempty,gt=0"`
}

type deleteDumpReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type coordinate struct {
	Latitude  *float64 `json:"latitude" validate:"required,latitude"`
	Longitude *float64 `json:"longitude" validate:"required,longitude"`
}
