package tempdumpentity

import (
	"capstonea03/be/src/libs/db/sql"
	applogger "capstonea03/be/src/libs/logger"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
)

type TempDumpModel struct {
	sql.Model
	MapSectorID *uuid.UUID  `gorm:"not null" json:"mapSectorId,omitempty"`
	Coordinate  *Coordinate `gorm:"not null" json:"coordinate,omitempty"`
	Type        *string     `gorm:"not null" json:"type,omitempty"`
	Capacity    *float64    `gorm:"not null" json:"capacity,omitempty"`
	Condition   *string     `gorm:"not null" json:"condition,omitempty"`
}

type Coordinate struct {
	Latitude  *float64 `json:"latitude" validate:"required,omitempty,latitude"`
	Longitude *float64 `json:"longitude" validate:"required,omitempty,longitude"`
}

func (TempDumpModel) TableName() string {
	return "temp_dumps"
}

func (c *Coordinate) Scan(value interface{}) error {
	return sonic.ConfigFastest.Unmarshal(value.([]byte), c)
}

func (c Coordinate) Value() (sql.Value, error) {
	return sonic.ConfigFastest.Marshal(c)
}

type tempDumpDB = sql.Service[TempDumpModel]

var tempDumpRepo *tempDumpDB
var logger = applogger.New("TempDumpModule")

func InitRepository(db *sql.DB) {
	if db == nil {
		logger.Panic("db cannot be nil")
	}

	tempDumpRepo = sql.NewService[TempDumpModel](db)
}

func TempDumpRepository() *tempDumpDB {
	if tempDumpRepo == nil {
		logger.Panic("tempDumpRepo is nil")
	}

	return tempDumpRepo
}
