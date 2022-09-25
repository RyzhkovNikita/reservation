package crud

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/lib/pq"
)

type BarProfile struct {
	Id           int    `orm:"auto"`
	Email        string `orm:"size(30)"`
	Name         string `orm:"size(50)"`
	Description  string `orm:"size(400)"`
	PasswordHash string `orm:"size(100)"`
	Address      string `orm:"size(100)"`
	LogoUrl      string
}

type BarCrud interface {
	GetById(id int) (*BarProfile, error)
	Insert(profile *BarProfile) (*BarProfile, error)
	AddToTransaction(profile *BarProfile) (*BarProfile, error)
}

func init() {
	err := orm.RegisterDataBase("database", "postgres", "")
	if err != nil {
		panic(err)
	}
}
