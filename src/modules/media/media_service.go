package media

import (
	"capstonea03/be/src/libs/db/mongo"
	"capstonea03/be/src/libs/localstorage"
	me "capstonea03/be/src/modules/media/media_entity"
)

func (*Module) getMediaService(id *mongo.ObjectID) (*me.MediaModel, error) {
	return me.Repository().FindOne(&mongo.FindOneOptions{
		Where: &[]mongo.Where{{
			{
				Key:   "_id",
				Value: id,
			},
		}},
	})
}

func (m *Module) addMediaService(data *me.MediaModel) (*me.MediaModel, error) {
	id, err := me.Repository().Create(data)
	if err != nil {
		return nil, err
	}

	return m.getMediaService(id)
}

func (*Module) deleteMediaService(id *mongo.ObjectID) error {
	return me.Repository().Destroy(&mongo.DestroyOptions{
		Where: &[]mongo.Where{{
			{
				Key:   "_id",
				Value: id,
			},
		}},
	})
}

func (m *Module) addMediaFileService(id *mongo.ObjectID, file *[]byte) error {
	return localstorage.SaveBinaryData(file, &localstorage.Option{
		Filename:       id.Hex(),
		Subdirectory:   "media",
		FilePermission: localstorage.OS_USER_RW,
	})
}
