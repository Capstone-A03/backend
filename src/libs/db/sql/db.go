package sql

import (
	"capstonea03/be/src/libs/gracefulshutdown"
	applogger "capstonea03/be/src/libs/logger"

	"gorm.io/gorm"
)

var logger = applogger.New("SQL")

func NewDB(dialector Dialector) *DB {
	logger.Log("initializing SQL database")

	db, err := gorm.Open(dialector, &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		logger.Panic(err)
	}

	gracefulshutdown.Add(gracefulshutdown.FnRunInShutdown{
		FnDescription: "closing SQL database",
		Fn: func() {
			db, err := db.DB()
			if err != nil {
				logger.Error(err)
				return
			}
			if err := db.Close(); err != nil {
				logger.Error(err)
			}
		},
	})

	return db
}
