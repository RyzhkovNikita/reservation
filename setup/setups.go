package setup

import (
	"barckend/storage"
	beego "github.com/beego/beego/v2/server/web"
)

type Setuper func() error

func addTestOwnerAccount() Setuper {
	return func() error {
		return nil
	}
}

func setImagesPath() Setuper {
	return func() error {
		beego.SetStaticPath("api/v1", storage.PathProvider.GetFilePath())
		return nil
	}
}

func addTestAdminAccount() Setuper {
	return func() error {
		return nil
	}
}

func GetSetupers() []Setuper {
	return []Setuper{
		addTestOwnerAccount(),
		addTestAdminAccount(),
		setImagesPath(),
	}
}
