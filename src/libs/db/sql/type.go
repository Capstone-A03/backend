package sql

import (
	"database/sql/driver"

	"gorm.io/gorm"
)

type DB = gorm.DB
type Dialector = gorm.Dialector
type Value = driver.Value
