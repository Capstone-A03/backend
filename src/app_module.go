package main

import (
	"strconv"
	"time"

	"capstonea03/be/src/libs/db/mongo"
	"capstonea03/be/src/libs/db/sql"
	"capstonea03/be/src/libs/db/sql/pg"
	"capstonea03/be/src/libs/env"
	"capstonea03/be/src/libs/hash/argon2"
	"capstonea03/be/src/libs/jwx/jwt"
	"capstonea03/be/src/libs/logger"
	"capstonea03/be/src/modules/auth"
	finaldump "capstonea03/be/src/modules/final_dump"
	logreport "capstonea03/be/src/modules/log_report"
	mapsector "capstonea03/be/src/modules/map_sector"
	"capstonea03/be/src/modules/mcu"
	tempdump "capstonea03/be/src/modules/temp_dump"
	"capstonea03/be/src/modules/truck"
	"capstonea03/be/src/modules/user"

	"github.com/gofiber/fiber/v2"
)

type module struct {
	app    *fiber.App
	logger *logger.Logger
}

func (m *module) load() {
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

	// JWT
	jwt.Init(&jwt.Config{
		Duration: func() *time.Duration {
			duration, err := time.ParseDuration(env.Get(env.JWT_DURATION))
			if err != nil {
				m.logger.Panic(err)
			}
			return &duration
		}(),
	})

	// Argon2
	argon2.Init(&argon2.Config{
		Memory: func() uint32 {
			hashMemory, err := strconv.ParseUint(env.Get(env.HASH_MEMORY), 10, 32)
			if err != nil {
				m.logger.Panic(err)
			}
			return uint32(hashMemory)
		}(),
		Iterations: func() uint32 {
			hashIterations, err := strconv.ParseUint(env.Get(env.HASH_ITERATIONS), 10, 32)
			if err != nil {
				m.logger.Panic(err)
			}
			return uint32(hashIterations)
		}(),
		Parallelism: func() uint8 {
			hashParallelism, err := strconv.ParseUint(env.Get(env.HASH_PARALLELISM), 10, 8)
			if err != nil {
				m.logger.Panic(err)
			}
			return uint8(hashParallelism)
		}(),
		SaltLength: func() int {
			hashSaltLength, err := strconv.Atoi(env.Get(env.HASH_SALTLENGTH))
			if err != nil {
				m.logger.Panic(err)
			}
			return hashSaltLength
		}(),
		KeyLength: func() uint32 {
			hashKeyLength, err := strconv.ParseUint(env.Get(env.HASH_KEYLENGTH), 10, 32)
			if err != nil {
				m.logger.Panic(err)
			}
			return uint32(hashKeyLength)
		}(),
	})

	m.controller()

	user.Load(&user.Module{
		App: m.app,
		DB:  pgDB,
	})

	auth.Load(&auth.Module{
		App: m.app,
	})

	mcu.Load(&mcu.Module{
		App: m.app,
		DB:  pgDB,
	})

	truck.Load(&truck.Module{
		App: m.app,
		DB:  pgDB,
	})

	mapsector.Load(&mapsector.Module{
		App: m.app,
		DB:  pgDB,
	})

	tempdump.Load(&tempdump.Module{
		App: m.app,
		DB:  pgDB,
	})

	finaldump.Load(&finaldump.Module{
		App: m.app,
		DB:  pgDB,
	})

	logreport.Load(&logreport.Module{
		App:      m.app,
		DBClient: mongoDBClient,
	})
}
