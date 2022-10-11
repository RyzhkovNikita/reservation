package crud

import (
	"github.com/beego/beego/v2/client/orm"
)

type Role uint

const (
	Owner Role = iota + 1
	Admin
	Guest
)

type PasswordHash struct {
	Id   uint64 `orm:"auto"`
	Hash string `orm:"size(256)"`
	User *User  `orm:"rel(one);on_delete(cascade)"`
}

type User struct {
	Id           uint64 `orm:"auto"`
	Role         Role
	IsActive     bool
	PasswordHash *PasswordHash `orm:"reverse(one)"`
	AdminInfo    *AdminInfo    `orm:"reverse(one)"`
	GuestInfo    *GuestInfo    `orm:"reverse(one)"`
}

type AdminInfo struct {
	Id         uint64 `orm:"auto"`
	Surname    string `orm:"size(50)"`
	Name       string `orm:"size(50)"`
	Patronymic string `orm:"size(50)"`
	Email      string `orm:"size(30)"`
	Phone      string `orm:"size(11)"`
	User       *User  `orm:"rel(one);on_delete(cascade)"`
}

type UpdateAdminInfo struct {
	Id         uint64
	Surname    *string
	Name       *string
	Patronymic *string
	Email      *string
	Phone      *string
}

type GuestInfo struct {
	Id    uint64 `orm:"auto"`
	Name  string `orm:"size(50)"`
	Phone string `orm:"size(11)"`
	User  *User  `orm:"rel(one);on_delete(cascade)"`
}

type Bar struct {
	Id                   uint64            `orm:"auto"`
	Email                string            `orm:"size(30)"`
	Address              string            `orm:"size(30)"`
	Name                 string            `orm:"size(30)"`
	Description          string            `orm:"size(30)"`
	LogoUrl              string            `orm:"size(30)"`
	Phone                string            `orm:"size(30)"`
	CreationTime         orm.DateTimeField `orm:"auto_now_add"`
	MaxReservTimeMinutes uint
	IsVisible            bool
}

type WorkHours struct {
	Id                   uint64 `orm:"auto"`
	From                 string `orm:"size(5)"`
	To                   string `orm:"size(5)"`
	MaxReservTimeMinutes uint
}

type Weekday struct {
	Id        uint   `orm:"auto"`
	Name      string `orm:"size(15)"`
	ShortName string `orm:"size(2)"`
}

type Table struct {
	Id          uint64 `orm:"auto"`
	Name        string `orm:"size(30)"`
	Description string `orm:"size(30)"`
	Email       string `orm:"size(30)"`
	Capacity    uint
	PhotoUrl    string `orm:"size(255)"`
	BarInfo     *Bar   `orm:"rel(one);on_delete(cascade)"`
}

type Reservation struct {
	Id           uint64            `orm:"auto"`
	From         string            `orm:"size(5)"`
	To           string            `orm:"size(5)"`
	CreationTime orm.DateTimeField `orm:"auto_now_add"`
	ModifyTime   orm.DateTimeField `orm:"auto_now"`
	Date         orm.DateField
	Table        *Table     `orm:"rel(one);on_delete(cascade)"`
	Guest        *GuestInfo `orm:"rel(one);on_delete(cascade)"`
}
