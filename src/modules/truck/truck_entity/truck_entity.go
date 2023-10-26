package truckentity

import (
	"capstonea03/be/src/libs/db/sql"
	applogger "capstonea03/be/src/libs/logger"

	"github.com/google/uuid"
)

type TruckModel struct {
	sql.Model
	MapSectorIDs    *[]*uuid.UUID `gorm:"column:map_sector_ids;type:uuid[]" json:"mapSectorIds,omitempty"`
	LicensePlate    *string       `gorm:"column:license_plate;uniqueIndex" json:"licensePlate,omitempty"`
	Type            *string       `gorm:"column:type;not null" json:"type,omitempty"`
	Capacity        *float64      `gorm:"column:capacity;not null" json:"capacity,omitempty"`
	FuelConsumption *float64      `gorm:"column:fuel_consumption" json:"fuelConsumption,omitempty"`
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

func Repository() *truckDB {
	if truckRepo == nil {
		logger.Panic("truckRepo is nil")
	}

	return truckRepo
}
