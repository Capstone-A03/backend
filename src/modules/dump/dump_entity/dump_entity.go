package dumpentity

import (
	"capstonea03/be/src/libs/db/sql"
	applogger "capstonea03/be/src/libs/logger"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
)

type DumpModel struct {
	sql.Model
	Name        *string     `gorm:"column:name;uniqueIndex;not null" json:"name,omitempty"`
	MapSectorID *uuid.UUID  `gorm:"column:map_sector_id" json:"mapSectorId,omitempty"`
	Coordinate  *Coordinate `gorm:"column:coordinate;not null" json:"coordinate,omitempty"`
	Type        *string     `gorm:"column:type;not null" json:"type,omitempty"`
	Capacity    *float64    `gorm:"column:capacity;not null" json:"capacity,omitempty"`
	Condition   *string     `gorm:"column:condition" json:"condition,omitempty"`
}

func (DumpModel) TableName() string {
	return "dumps"
}

type DumpType string

const (
	TempDump  DumpType = "TEMP_DUMP"
	FinalDump DumpType = "FINAL_DUMP"
)

type Coordinate struct {
	Latitude  *float64 `json:"latitude" validate:"omitempty,latitude"`
	Longitude *float64 `json:"longitude" validate:"omitempty,longitude"`
}

func (c *Coordinate) Scan(value interface{}) error {
	return sonic.ConfigFastest.Unmarshal(value.([]byte), c)
}

func (c Coordinate) Value() (sql.Value, error) {
	return sonic.ConfigFastest.Marshal(c)
}

type dumpDB = sql.Service[DumpModel]

var dumpRepo *dumpDB
var logger = applogger.New("DumpModule")

func InitRepository(db *sql.DB) {
	if db == nil {
		logger.Panic("db cannot be nil")
	}

	dumpRepo = sql.NewService[DumpModel](db)
}

func Repository() *dumpDB {
	if dumpRepo == nil {
		logger.Panic("dumpRepo is nil")
	}

	return dumpRepo
}
