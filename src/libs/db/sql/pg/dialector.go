package pg

import (
	"capstonea03/be/src/libs/logger"
	"capstonea03/be/src/libs/validator"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Address      string `validate:"gt=0"`
	User         string `validate:"gt=0"`
	Password     string `validate:"gt=0"`
	DatabaseName string `validate:"gt=0"`
}

func NewDialector(config *Config) gorm.Dialector {
	logger := logger.New("PgDialector")

	if err := validator.Struct(config); err != nil {
		logger.Panic(err)
	}

	return postgres.New(postgres.Config{
		DSN: "postgresql://" + config.User + ":" + config.Password + "@" + config.Address + "/" + config.DatabaseName,
	})
}
