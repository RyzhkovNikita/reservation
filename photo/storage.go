package photo

import (
	"barckend/storage"
	"github.com/pkg/errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Storage interface {
	SavePhoto(info Input) (uint64, error)
	GetPhoto(id uint64) (Info, error)
}

type Input struct {
	ByteReader io.Reader
	Ext        storage.ImageExtension
}

type Info struct {
	Id       uint64
	FileName string
	File     *os.File
}

var NoImageError = errors.New("No image with provided id")

func GetPhotoStorage() Storage {
	return instance
}

var instance = &photoStorageImpl{}

type photoStorageImpl struct{}

func (p photoStorageImpl) SavePhoto(info Input) (uint64, error) {
	id := uint64(time.Now().UnixMilli())
	create, err := os.Create(
		storage.PathProvider.GetImagePath(
			strconv.FormatUint(id, 10) + "." + info.Ext.ToString(),
		),
	)
	if err != nil {
		return 0, errors.Wrap(err, "Error when creating file")
	}
	defer create.Close()
	_, err = io.Copy(create, info.ByteReader)
	if err != nil {
		return 0, errors.New("Error while copying file")
	}
	return id, nil
}

func (p photoStorageImpl) GetPhoto(id uint64) (Info, error) {
	imageName := strconv.FormatUint(id, 10)
	entries, err := os.ReadDir(storage.PathProvider.GetImagesDirPath())
	if err != nil {
		return Info{}, errors.Wrap(err, "Error when reading directory")
	}
	if len(entries) == 0 {
		return Info{}, NoImageError
	}
	var needEntry *os.DirEntry
	for _, entry := range entries {
		if !strings.Contains(entry.Name(), imageName) {
			continue
		}
		needEntry = &entry
		break
	}
	if needEntry == nil {
		return Info{}, NoImageError
	}
	file, err := os.Open(storage.PathProvider.GetImagePath((*needEntry).Name()))
	return Info{
		Id:   id,
		File: file,
	}, nil
}
