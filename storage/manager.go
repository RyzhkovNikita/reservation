package storage

import (
	"barckend/conf"
	"os"
)

var PathProvider IPathProvider = &pathProviderImpl{}

type IPathProvider interface {
	GetFilePath() string
	GetImagesDirPath() string
	GetImagePath(imageName string) string
}

type pathProviderImpl struct{}

func (p *pathProviderImpl) GetFilePath() string {
	return conf.AppConfig.FileStoragePath
}

func (p *pathProviderImpl) GetImagesDirPath() string {
	return conf.AppConfig.FileStoragePath + string(os.PathSeparator) + "images"
}

func (p *pathProviderImpl) GetImagePath(imageName string) string {
	return p.GetImagesDirPath() + string(os.PathSeparator) + imageName
}
