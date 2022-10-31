package crud

import (
	"barckend/conf"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
)

var NoRowsUpdated = fmt.Errorf("no rows has been updated")

func init() {
	orm.Debug = conf.AppConfig.Mode != conf.Prod
	err := orm.RegisterDataBase(
		"default",
		conf.AppConfig.DriverName,
		conf.AppConfig.DataSourceUrl,
	)
	if err != nil {
		panic(err)
	}
	orm.RegisterModel(
		new(PasswordHash),
		new(User),
		new(AdminInfo),
		new(OwnerInfo),
		new(GuestInfo),
		new(Bar),
		new(WorkHours),
		new(Weekday),
		new(Table),
		new(Reservation),
	)
	name := "default"
	force := false
	verbose := true
	err = orm.RunSyncdb(name, force, verbose)
	if err != nil {
		panic(fmt.Errorf("error occured while Syncdb: %s", err))
	}
	newOrm := orm.NewOrm()
	Db = &userCrudImpl{
		ormer: newOrm,
	}
	initializeBarCrud(newOrm)
}
