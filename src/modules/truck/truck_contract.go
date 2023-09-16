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
	LicensePlate    *string  `json:"licensePlate" validate:"required,gt=0"`
	Type            *string  `json:"type" validate:"required,gt=0"`
	Capacity        *float64 `json:"capacity" validate:"required,gt=0"`
	FuelConsumption *float64 `json:"fuelConsumption" validate:"required,gt=0"`
}

type updateTruckReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type updateTruckReq struct {
	LicensePlate    *string  `json:"licensePlate" validate:"omitempty,gt=0"`
	Type            *string  `json:"type" validate:"omitempty,gt=0"`
	Capacity        *float64 `json:"capacity" validate:"omitempty,gt=0"`
	FuelConsumption *float64 `json:"fuelConsumption" validate:"omitempty,gt=0"`
}

type deleteTruckReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}
