package env

import (
	"fmt"
	"os"

	applogger "capstonea03/be/src/libs/logger"

	_ "github.com/joho/godotenv/autoload"
)

type Env string

const (
	APP_NAME             Env = "APP_NAME"
	APP_MODE             Env = "APP_MODE"
	APP_ADDRESS          Env = "APP_ADDRESS"
	APP_PUBLIC_DIRECTORY Env = "APP_PUBLIC_DIRECTORY"

	WEB_ADDRESS Env = "WEB_ADDRESS"

	POSTGRES_ADDRESS  Env = "POSTGRES_ADDRESS"
	POSTGRES_USER     Env = "POSTGRES_USER"
	POSTGRES_PASSWORD Env = "POSTGRES_PASSWORD"
	POSTGRES_DB       Env = "POSTGRES_DB"

	MONGO_ADDRESS              Env = "MONGO_ADDRESS"
	MONGO_INITDB_ROOT_USERNAME Env = "MONGO_INITDB_ROOT_USERNAME"
	MONGO_INITDB_ROOT_PASSWORD Env = "MONGO_INITDB_ROOT_PASSWORD"
	MONGO_DATABASE_NAME        Env = "MONGO_DATABASE_NAME"

	JWT_DURATION Env = "JWT_DURATION"

	HASH_MEMORY      Env = "HASH_MEMORY"
	HASH_ITERATIONS  Env = "HASH_ITERATIONS"
	HASH_PARALLELISM Env = "HASH_PARALLELISM"
	HASH_SALTLENGTH  Env = "HASH_SALTLENGTH"
	HASH_KEYLENGTH   Env = "HASH_KEYLENGTH"

	INITIAL_USER_NAME     Env = "INITIAL_USER_NAME"
	INITIAL_USER_USERNAME Env = "INITIAL_USER_USERNAME"
	INITIAL_USER_PASSWORD Env = "INITIAL_USER_PASSWORD"
	INITIAL_USER_ROLE     Env = "INITIAL_USER_ROLE"
)

var logger = applogger.New("Env")

type Option struct {
	MustExist bool
}

func Get(env Env, option ...Option) string {
	envVal, isExist := os.LookupEnv(string(env))
	if !isExist {
		err := fmt.Errorf("unknown env variable: %s", env)
		if len(option) > 0 && option[0].MustExist {
			logger.Panic(err)
		}
		logger.Error(err)
	}

	return envVal
}
