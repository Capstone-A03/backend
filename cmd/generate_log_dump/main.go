package main

import (
	"capstonea03/be/src/libs/db/mongo"
	"capstonea03/be/src/libs/db/sql"
	"capstonea03/be/src/libs/db/sql/pg"
	"capstonea03/be/src/libs/env"
	applogger "capstonea03/be/src/libs/logger"
	de "capstonea03/be/src/modules/dump/dump_entity"
	lde "capstonea03/be/src/modules/log_dump/log_dump_entity"
	"capstonea03/be/src/utils"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	logger := applogger.New("Generate Log Dump")

	// PostgreSQL database
	pgDB := sql.NewDB(pg.NewDialector(&pg.Config{
		Address:      env.Get(env.POSTGRES_ADDRESS),
		User:         env.Get(env.POSTGRES_USER),
		Password:     env.Get(env.POSTGRES_PASSWORD),
		DatabaseName: env.Get(env.POSTGRES_DB),
	}))

	// MongoDB database
	mongoDBClient := mongo.NewClient(&mongo.Config{
		Address:  env.Get(env.MONGO_ADDRESS),
		User:     env.Get(env.MONGO_INITDB_ROOT_USERNAME),
		Password: env.Get(env.MONGO_INITDB_ROOT_PASSWORD),
	})

	de.InitRepository(pgDB)
	lde.InitRepository(mongoDBClient)

	dumpListData, page, err := de.Repository().FindAll(&sql.FindAllOptions{
		Where: &[]sql.FindAllWhere{{
			Where: sql.Where{
				Query: "type = ?",
				Args:  []interface{}{de.TempDump},
			},
		}},
		Limit: utils.AsRef(-1),
	})
	if err != nil {
		logger.Panic(err.Error())
	}
	if page.Count == 0 {
		logger.Panic("No dump data")
	}

	createLogDumpListChs := make([]chan bool, len(*dumpListData))
	for i := range createLogDumpListChs {
		createLogDumpListChs[i] = make(chan bool)
	}
	for i, dumpData := range *dumpListData {
		dumpID := dumpData.ID
		idx := i
		dumpCapacity := *dumpData.Capacity
		go func() {
			randomVolume := rand.Float64() * dumpCapacity
			logDumpID, err := lde.Repository().Create(&lde.LogDumpModel{
				DumpID:         dumpID,
				MeasuredVolume: &randomVolume,
			})
			if err != nil {
				logger.Error(err.Error())
				createLogDumpListChs[idx] <- false
				return
			}
			logger.Log(fmt.Sprintf("ID: %s | DumpID: %s | Volume: %f", logDumpID, dumpID, randomVolume))
			createLogDumpListChs[idx] <- true
		}()
	}

	successCount := 0
	failCount := 0
	for i := range createLogDumpListChs {
		if isSuccess := <-createLogDumpListChs[i]; isSuccess {
			successCount++
		} else {
			failCount++
		}
		close(createLogDumpListChs[i])
	}
	logger.Log(fmt.Sprintf("Success: %d | Fail: %d", successCount, failCount))
}
