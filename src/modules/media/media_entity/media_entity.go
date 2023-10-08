package mediaentity

import (
	"capstonea03/be/src/libs/db/mongo"
	"capstonea03/be/src/libs/env"
	applogger "capstonea03/be/src/libs/logger"
)

type MediaModel struct {
	mongo.Model `bson:"inline"`
	Filename    *string `bson:"filename,omitempty" json:"filename,omitempty"`
	Header      *string `bson:"header,omitempty" json:"header,omitempty"`
	Size        *int    `bson:"size,omitempty" json:"size,omitempty"`
}

func (MediaModel) DatabaseName() string {
	return env.Get(env.MONGO_DATABASE_NAME)
}

func (MediaModel) CollectionName() string {
	return "media"
}

type mediaDB = mongo.Service[MediaModel]

var mediaRepo *mediaDB
var logger = applogger.New("MediaModule")

func InitRepository(client *mongo.Client) {
	if client == nil {
		logger.Panic("client cannot be nil")
	}

	mediaRepo = mongo.NewService[MediaModel](client)
}

func Repository() *mediaDB {
	if mediaRepo == nil {
		logger.Panic("mediaRepo is nil")
	}

	return mediaRepo
}
