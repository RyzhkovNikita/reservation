package setup

import (
	"barckend/storage"
	beego "github.com/beego/beego/v2/server/web"
	"os"
)

type Setuper func() error

func addTestOwnerAccount() Setuper {
	return func() error {
		return nil
	}
}

func setImagesPath() Setuper {
	return func() error {
		beego.SetStaticPath("api/v1/files", storage.PathProvider.GetFilePath())
		beego.SetStaticPath("api/v1/images", storage.PathProvider.GetImagesDirPath())
		return nil
	}
}

func addTestAdminAccount() Setuper {
	return func() error {
		return nil
	}
}

func createFilePath() Setuper {
	return func() error {
		if _, err := os.Stat(storage.PathProvider.GetFilePath()); os.IsNotExist(err) {
			err = os.Mkdir(storage.PathProvider.GetFilePath(), os.FileMode(0700))
			if err != nil {
				return err
			}
		}
		if _, err := os.Stat(storage.PathProvider.GetImagesDirPath()); os.IsNotExist(err) {
			err = os.Mkdir(storage.PathProvider.GetImagesDirPath(), os.FileMode(0700))
			return err
		}
		return nil
	}
}

func GetSetupers() []Setuper {
	return []Setuper{
		addTestOwnerAccount(),
		addTestAdminAccount(),
		setImagesPath(),
		createFilePath(),
	}
}
