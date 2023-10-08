package mcuentity

import (
	"capstonea03/be/src/libs/db/sql"
	applogger "capstonea03/be/src/libs/logger"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
)

type McuModel struct {
	sql.Model
	DumpID     *uuid.UUID  `gorm:"column:dump_id;not null" json:"dumpId,omitempty"`
	Coordinate *Coordinate `gorm:"column:coordinate;not null" json:"coordinate,omitempty"`
}

func (McuModel) TableName() string {
	return "mcus"
}

type Coordinate struct {
	Latitude  *float64 `json:"latitude" validate:"required,omitempty,latitude"`
	Longitude *float64 `json:"longitude" validate:"required,omitempty,longitude"`
}

func (c *Coordinate) Scan(value interface{}) error {
	return sonic.ConfigFastest.Unmarshal(value.([]byte), c)
}

func (c Coordinate) Value() (sql.Value, error) {
	return sonic.ConfigFastest.Marshal(c)
}

type mcuDB = sql.Service[McuModel]

var mcuRepo *mcuDB
var logger = applogger.New("McuModule")

func InitRepository(db *sql.DB) {
	if db == nil {
		logger.Panic("db cannot be nil")
	}

	mcuRepo = sql.NewService[McuModel](db)
}

func Repository() *mcuDB {
	if mcuRepo == nil {
		logger.Panic("mcuRepo is nil")
	}

	return mcuRepo
}
