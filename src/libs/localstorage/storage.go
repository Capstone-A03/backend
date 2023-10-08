package localstorage

import (
	applogger "capstonea03/be/src/libs/logger"
	"capstonea03/be/src/libs/validator"
	"io/fs"
	"os"
)

type Config struct {
	RootDirectory string `validate:"required,startswith=/,endsnotwith=/"`
}

var storageConfig *Config
var logger = applogger.New("LocalStorage")

func New(config *Config) {
	logger.Log("initializing local storage")

	if err := validator.Struct(config); err != nil {
		logger.Panic(err)
	}

	storageConfig = config
}

type Option struct {
	Filename       string `validate:"required,excludes=/"`
	Subdirectory   string `validate:"startsnotwith=/,endsnotwith=/"`
	FilePermission fs.FileMode
}

func SaveBinaryData(data *[]byte, option *Option) error {
	if err := validator.Struct(option); err != nil {
		logger.Error(err)
		return err
	}

	path := storageConfig.RootDirectory
	if len(option.Subdirectory) > 0 {
		path += "/" + option.Subdirectory
	}
	if err := os.MkdirAll(path, os.ModeDir|OS_USER_RWX|OS_GROUP_RX|OS_ALL_RX); err != nil {
		logger.Error(err)
		return err
	}

	path += "/" + option.Filename

	permission := OS_USER_R
	if option.FilePermission > 0 {
		permission = option.FilePermission
	}

	if err := os.WriteFile(path, *data, permission); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func RemoveFile(option *Option) error {
	if err := validator.Struct(option); err != nil {
		logger.Error(err)
		return err
	}

	path := storageConfig.RootDirectory
	if len(option.Subdirectory) > 0 {
		path += "/" + option.Subdirectory
	}
	path += "/" + option.Filename

	if err := os.Remove(path); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
