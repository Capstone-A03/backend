package truck

import "github.com/google/uuid"

type getTruckListReqQuery struct {
	LastID *uuid.UUID `query:"lastId"`
	Limit  *int       `query:"limit"`
}

type getTruckReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type addTruckReq struct {
	LicensePlate    *string `json:"licensePlate" validate:"required,gt=0"`
	Type            *string `json:"type" validate:"required,gt=0"`
	Capacity        *int    `json:"capacity" validate:"required,gt=0"`
	FuelConsumption *int    `json:"fuelConsumption" validate:"required,gt=0"`
}

type updateTruckReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type updateTruckReq struct {
	LicensePlate    *string `json:"licensePlate"`
	Type            *string `json:"type"`
	Capacity        *int    `json:"capacity"`
	FuelConsumption *int    `json:"fuelConsumption"`
}

type deleteTruckReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}
