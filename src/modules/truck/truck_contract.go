package truck

import "github.com/google/uuid"

type getTruckListReqQuery struct {
	SearchByIsActive *bool      `query:"isActive"`
	LastID           *uuid.UUID `query:"lastId"`
	Limit            *int       `query:"limit"`
}

type getTruckReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type addTruckReq struct {
	MapSectorIDs    *[]*uuid.UUID `json:"mapSectorIds" validate:"omitempty,unique"`
	LicensePlate    *string       `json:"licensePlate" validate:"omitempty,gt=0"`
	Type            *string       `json:"type" validate:"gt=0"`
	Capacity        *float64      `json:"capacity" validate:"gt=0"`
	FuelConsumption *float64      `json:"fuelConsumption" validate:"omitempty,gt=0"`
	IsActive        *bool         `json:"isActive" validate:"omitempty"`
}

type updateTruckReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type updateTruckReq struct {
	MapSectorIDs    *[]*uuid.UUID `json:"mapSectorIds" validate:"omitempty,unique"`
	LicensePlate    *string       `json:"licensePlate" validate:"omitempty,gt=0"`
	Type            *string       `json:"type" validate:"omitempty,gt=0"`
	Capacity        *float64      `json:"capacity" validate:"omitempty,gt=0"`
	FuelConsumption *float64      `json:"fuelConsumption" validate:"omitempty,gt=0"`
	IsActive        *bool         `json:"isActive" validate:"omitempty"`
}

type deleteTruckReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}
