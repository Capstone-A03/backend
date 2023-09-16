package finaldumpentity

import (
	"capstonea03/be/src/libs/db/sql"
	applogger "capstonea03/be/src/libs/logger"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
)

type FinalDumpModel struct {
	sql.Model
	Name        *string     `gorm:"uniqueIndex;not null" json:"name,omitempty"`
	MapSectorID *uuid.UUID  `gorm:"not null" json:"mapSectorId,omitempty"`
	Coordinate  *Coordinate `gorm:"not null" json:"coordinate,omitempty"`
}

type Coordinate struct {
	Latitude  *float64 `json:"latitude" validate:"required,omitempty,latitude"`
	Longitude *float64 `json:"longitude" validate:"required,omitempty,longitude"`
}

func (FinalDumpModel) TableName() string {
	return "final_dumps"
}

func (c *Coordinate) Scan(value interface{}) error {
	return sonic.ConfigFastest.Unmarshal(value.([]byte), c)
}

func (c Coordinate) Value() (sql.Value, error) {
	return sonic.ConfigFastest.Marshal(c)
}

type finalDumpDB = sql.Service[FinalDumpModel]

var finalDumpRepo *finalDumpDB
var logger = applogger.New("FinalDumpModule")

func InitRepository(db *sql.DB) {
	if db == nil {
		logger.Panic("db cannot be nil")
	}

	finalDumpRepo = sql.NewService[FinalDumpModel](db)
}

func FinalDumpRepository() *finalDumpDB {
	if finalDumpRepo == nil {
		logger.Panic("finalDumpRepo is nil")
	}

	return finalDumpRepo
}
