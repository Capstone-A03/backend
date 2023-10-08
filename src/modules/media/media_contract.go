package media

import (
	"capstonea03/be/src/libs/db/mongo"
	"mime/multipart"
)

type getMediaReqParam struct {
	ID *mongo.ObjectID `params:"id" validate:"required"`
}

type addMediaReqForm struct {
	File *multipart.FileHeader `form:"file" validate:"required"`
}
