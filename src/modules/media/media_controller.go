package media

import (
	"capstonea03/be/src/contracts"
	"capstonea03/be/src/libs/db/mongo"
	"capstonea03/be/src/libs/parser"
	me "capstonea03/be/src/modules/media/media_entity"
	"capstonea03/be/src/utils"

	"github.com/gofiber/fiber/v2"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/media/:id", m.getMedia)
	m.App.Post("/api/v1/media", m.addMedia)
}

func (m *Module) getMedia(c *fiber.Ctx) error {
	param := new(getMediaReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	mediaData, err := m.getMediaService(param.ID)
	if err != nil {
		if mongo.IsErrNoDocuments(err) {
			return contracts.NewError(fiber.ErrNotFound, err.Error())
		}
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: mediaData,
	})
}

func (m *Module) addMedia(c *fiber.Ctx) error {
	form := new(addMediaReqForm)
	if err := parser.ParseReqMultipartForm(c, form); err != nil {
		return contracts.NewError(fiber.ErrBadRequest, err.Error())
	}

	file, header, err := parser.ParseMultipartFileToBytes(form.File)
	if err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	mediaData, err := m.addMediaService(&me.MediaModel{
		Filename: &form.File.Filename,
		Header:   utils.AsRef(string(*header)),
		Size:     utils.AsRef(len(*file)),
	})
	if err != nil {
		return contracts.NewError(fiber.ErrInternalServerError, err.Error())
	}

	if err := m.addMediaFileService(mediaData.ID, file); err != nil {
		errStr := err.Error()
		if err := m.deleteMediaService(mediaData.ID); err != nil {
			errStr = err.Error()
		}
		return contracts.NewError(fiber.ErrInternalServerError, errStr)
	}

	return c.Status(fiber.StatusOK).JSON(&contracts.Response{
		Data: mediaData,
	})
}
