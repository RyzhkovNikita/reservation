package storage

import (
	"barckend/conf"
	"os"
)

var PathProvider = &pathProviderImpl{}

type IPathProvider interface {
	GetFilePath() string
	GetImagePath() string
}

type pathProviderImpl struct{}

func (p *pathProviderImpl) GetFilePath() string {
	return conf.AppConfig.FileStoragePath
}

func (p *pathProviderImpl) GetImagePath() string {
	return conf.AppConfig.FileStoragePath + string(os.PathSeparator) + "images"
}
