package mongo

import (
	"context"

	"capstonea03/be/src/libs/gracefulshutdown"
	applogger "capstonea03/be/src/libs/logger"
	"capstonea03/be/src/libs/validator"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Address  string `validate:"gt=0"`
	User     string `validate:"gt=0"`
	Password string `validate:"gt=0"`
}

var logger = applogger.New("MongoDB")

func NewClient(config *Config) *Client {
	logger.Log("initializing MongoDB client")

	if err := validator.Struct(config); err != nil {
		logger.Panic(err)
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://"+config.User+":"+config.Password+"@"+config.Address))
	if err != nil {
		logger.Panic(err)
	}

	gracefulshutdown.Add(gracefulshutdown.FnRunInShutdown{
		FnDescription: "disconnecting mongodb client",
		Fn: func() {
			if err := client.Disconnect(context.TODO()); err != nil {
				logger.Error(err)
			}
		},
	})

	return client
}
