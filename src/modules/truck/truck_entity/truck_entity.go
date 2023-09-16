package truckentity

import (
	"capstonea03/be/src/libs/db/sql"
	applogger "capstonea03/be/src/libs/logger"
)

type TruckModel struct {
	sql.Model
	LicensePlate    *string `gorm:"uniqueIndex;not null" json:"licensePlate,omitempty"`
	Type            *string `gorm:"not null" json:"type,omitempty"`
	Capacity        *int    `gorm:"not null" json:"capacity,omitempty"`
	FuelConsumption *int    `gorm:"not null" json:"fuelConsumption,omitempty"`
}

func (TruckModel) TableName() string {
	return "trucks"
}

type truckDB = sql.Service[TruckModel]

var truckRepo *truckDB
var logger = applogger.New("TruckModule")

func InitRepository(db *sql.DB) {
	if db == nil {
		logger.Panic("db cannot be nil")
	}

	truckRepo = sql.NewService[TruckModel](db)
}

func TruckRepository() *truckDB {
	if truckRepo == nil {
		logger.Panic("truckRepo is nil")
	}

	return truckRepo
}
