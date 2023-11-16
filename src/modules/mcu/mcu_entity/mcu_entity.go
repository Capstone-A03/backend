package mcuentity

import (
	"capstonea03/be/src/libs/db/sql"
	applogger "capstonea03/be/src/libs/logger"

	"github.com/google/uuid"
)

type McuModel struct {
	sql.Model
	DumpID *uuid.UUID `gorm:"column:dump_id;not null" json:"dumpId,omitempty"`
}

func (McuModel) TableName() string {
	return "mcus"
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
