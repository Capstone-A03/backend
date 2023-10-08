package mapsectorentity

import (
	"capstonea03/be/src/libs/db/sql"
	applogger "capstonea03/be/src/libs/logger"

	"github.com/bytedance/sonic"
)

type MapSectorModel struct {
	sql.Model
	Name    *string      `gorm:"column:name;uniqueIndex;not null" json:"name,omitempty"`
	Polygon *Coordinates `gorm:"column:polygon;not null" json:"polygon,omitempty"`
}

func (MapSectorModel) TableName() string {
	return "map_sectors"
}

type Coordinates []*Coordinate

func (c *Coordinates) Scan(value interface{}) error {
	return sonic.ConfigFastest.Unmarshal(value.([]byte), c)
}

func (c Coordinates) Value() (sql.Value, error) {
	return sonic.ConfigFastest.Marshal(c)
}

type Coordinate struct {
	Latitude  *float64 `json:"latitude" validate:"required,omitempty,latitude"`
	Longitude *float64 `json:"longitude" validate:"required,omitempty,longitude"`
}

type mapSectorDB = sql.Service[MapSectorModel]

var mapSectorRepo *mapSectorDB
var logger = applogger.New("MapSectorModule")

func InitRepository(db *sql.DB) {
	if db == nil {
		logger.Panic("db cannot be nil")
	}

	mapSectorRepo = sql.NewService[MapSectorModel](db)
}

func Repository() *mapSectorDB {
	if mapSectorRepo == nil {
		logger.Panic("mapSectorRepo is nil")
	}

	return mapSectorRepo
}
