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
	AdminInfo    *AdminInfo    `orm:"null;reverse(one)"`
	GuestInfo    *GuestInfo    `orm:"null;reverse(one)"`
	OwnerInfo    *OwnerInfo    `orm:"null;reverse(one)"`
}

func (user *User) IsAdmin() bool {
	return user.Role == Admin
}

func (user *User) IsOwner() bool {
	return user.Role == Owner
}

type AdminInfo struct {
	Id         uint64 `orm:"auto"`
	Surname    string `orm:"size(50)"`
	Name       string `orm:"size(50)"`
	Patronymic string `orm:"size(50)"`
	Email      string `orm:"size(30)"`
	Phone      string `orm:"size(11)"`
	User       *User  `orm:"rel(one);on_delete(cascade)"`
	Bar        *Bar   `orm:"null;rel(fk);on_delete(set_null)"`
}

type OwnerInfo struct {
	Id         uint64 `orm:"auto"`
	Surname    string `orm:"size(50)"`
	Name       string `orm:"size(50)"`
	Patronymic string `orm:"size(50)"`
	Email      string `orm:"size(30)"`
	Phone      string `orm:"size(11)"`
	User       *User  `orm:"rel(one);on_delete(cascade)"`
	Bars       []*Bar `orm:"reverse(many)"`
}

type GuestInfo struct {
	Id    uint64 `orm:"auto"`
	Name  string `orm:"size(50)"`
	Phone string `orm:"size(11)"`
	User  *User  `orm:"rel(one);on_delete(cascade)"`
}

type Bar struct {
	Id           uint64            `orm:"auto"`
	Email        string            `orm:"size(30)"`
	Address      string            `orm:"size(30)"`
	Name         string            `orm:"size(30)"`
	Description  string            `orm:"size(30)"`
	LogoUrl      string            `orm:"size(30)"`
	Phone        string            `orm:"size(30)"`
	CreationTime orm.DateTimeField `orm:"auto_now_add"`
	IsVisible    bool
	OwnerInfo    *OwnerInfo   `orm:"null;rel(fk);on_delete(set_null)"`
	Admins       []*AdminInfo `orm:"reverse(many)"`
	WorkHours    []*WorkHours `orm:"reverse(many)"`
}

type WorkHours struct {
	Id      uint64 `orm:"auto"`
	Weekday uint
	From    orm.DateTimeField
	To      orm.DateTimeField
	Bar     *Bar `orm:"rel(fk);on_delete(cascade)"`
}

type Weekday struct {
	Id        uint   `orm:"auto"`
	Name      string `orm:"size(15)"`
	ShortName string `orm:"size(2)"`
}

type Table struct {
	Id          uint64 `orm:"auto"`
	Name        string `orm:"size(30)"`
	Description string `orm:"size(100)"`
	Capacity    uint
	BarInfo     *Bar `orm:"rel(fk);on_delete(cascade)"`
}

type Reservation struct {
	Id           uint64 `orm:"auto"`
	From         orm.DateTimeField
	To           orm.DateTimeField
	Date         orm.DateField
	PersonCount  uint8
	Comment      *string           `orm:"null;size(100)"`
	CreationTime orm.DateTimeField `orm:"auto_now_add"`
	ModifyTime   orm.DateTimeField `orm:"auto_now"`
	Table        *Table            `orm:"rel(fk);on_delete(cascade)"`
	Guest        *GuestInfo        `orm:"null;rel(fk);on_delete(cascade)"`
}

func (r Reservation) IsInterceptingWith(other Reservation) bool {
	thisFrom, thisTo := r.From.Value(), r.To.Value()
	otherFrom, otherTo := other.From.Value(), other.To.Value()
	if thisFrom.Equal(otherFrom) || thisTo.Equal(otherTo) {
		return true
	}
	if thisTo.Before(otherFrom) ||
		thisTo.Equal(otherFrom) ||
		thisFrom.After(otherTo) ||
		thisFrom.Equal(otherTo) {
		return false
	}
	return true
}
