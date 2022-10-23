package link

import (
	"barckend/photo"
	"barckend/storage"
	"fmt"
	"github.com/pkg/errors"
)

type Builder interface {
	GetLinkForFile(host string, fileInfo storage.FileInfo) (string, error)
	GetLinkForPhoto(host string, photoInfo photo.Info) (string, error)
}

var instance = &builderImpl{}

type builderImpl struct{}

func GetLinkProvider() Builder {
	return instance
}

func (b builderImpl) GetLinkForFile(host string, fileInfo storage.FileInfo) (string, error) {
	return fmt.Sprintf("%s%s", host, fileInfo.Id), nil
}

func (b builderImpl) GetLinkForPhoto(host string, photoInfo photo.Info) (string, error) {
	fileInfo, err := photoInfo.File.Stat()
	if err != nil {
		return "", errors.Wrap(err, "Error when get photo name")
	}
	return fmt.Sprintf("%s/api/v1/images/%s", host, fileInfo.Name()), nil
}
