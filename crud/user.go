package crud

import (
	"barckend/conf"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/pkg/errors"

	//_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var Db AdminCrud

var NoRowsUpdated = fmt.Errorf("no rows has been updated")

func init() {
	b := conf.AppConfig.Mode != conf.Prod
	orm.Debug = b
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
		new(GuestInfo),
		new(Bar),
		new(WorkHours),
		new(Weekday),
		new(Table),
		new(Reservation),
	)
	name := "default"
	force := true
	verbose := true
	err = orm.RunSyncdb(name, force, verbose)
	if err != nil {
		panic(fmt.Errorf("error occured while Syncdb: %s", err))
	}
	Db = &barCrudImpl{
		ormer: orm.NewOrm(),
	}
}

type AdminCrud interface {
	GetById(id uint64) (*AdminInfo, error)
	IsEmailOccupied(email string, ownerId *uint64) (bool, error)
	IsPhoneOccupied(phone string, ownerId *uint64) (bool, error)
	Insert(profile *AdminInfo) (*AdminInfo, error)
	Update(profile *UpdateAdminInfo) (*AdminInfo, error)
	CheckCredentialsEmail(email string, passwordHash string) (*AdminInfo, error)
	CheckCredentialsPhone(phone string, passwordHash string) (*AdminInfo, error)
}

type barCrudImpl struct {
	ormer orm.Ormer
}

func (b *barCrudImpl) IsEmailOccupied(email string, ownerId *uint64) (bool, error) {
	adminInfo := &AdminInfo{}
	expr := b.ormer.QueryTable(adminInfo).Filter("Email", email)
	if ownerId != nil {
		expr = expr.Exclude("id", ownerId)
	}
	err := expr.
		One(adminInfo)
	if err == orm.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, errors.Wrap(err, fmt.Sprintf("Unknown error"))
	}
	return true, nil
}

func (b *barCrudImpl) IsPhoneOccupied(phone string, ownerId *uint64) (bool, error) {
	adminInfo := &AdminInfo{}
	expr := b.ormer.QueryTable(adminInfo).Filter("Phone", phone)
	if ownerId != nil {
		expr = expr.Exclude("id", ownerId)
	}
	err := expr.One(adminInfo)
	if err == orm.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, errors.Wrap(err, fmt.Sprintf("Unknown error"))
	}
	return true, nil
}

func (b *barCrudImpl) Insert(profile *AdminInfo) (*AdminInfo, error) {
	_, err := b.ormer.Insert(profile.User)
	if err != nil {
		return nil, errors.Wrap(
			err,
			fmt.Sprintf("Insert user error occured. User: %v", profile.User),
		)
	}
	_, err = b.ormer.Insert(profile.User.PasswordHash)
	if err != nil {
		return nil, errors.Wrap(
			err,
			fmt.Sprintf("Insert password error occured. Password hash: %v", profile.User.PasswordHash),
		)
	}
	_, err = b.ormer.Insert(profile)
	if err != nil {
		return nil, errors.Wrap(
			err,
			fmt.Sprintf("Insert admin info error occured. AdminInfo: %v", profile),
		)
	}
	return profile, nil
}

func (b *barCrudImpl) CheckCredentialsEmail(email string, passwordHash string) (*AdminInfo, error) {
	user := &User{}
	err := b.ormer.
		QueryTable(user).
		Filter("AdminInfo__email", email).
		Filter("PasswordHash__hash", passwordHash).
		RelatedSel().
		One(user)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Unknown error for email: %s", email))
	}
	_, err = b.ormer.LoadRelated(user, "AdminInfo")
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Unknown error when loading related for email: %s", email))
	}
	user.AdminInfo.User = user
	return user.AdminInfo, nil
}

func (b *barCrudImpl) CheckCredentialsPhone(phone string, passwordHash string) (*AdminInfo, error) {
	user := &User{}
	err := b.ormer.
		QueryTable(user).
		Filter("AdminInfo__phone", phone).
		Filter("PasswordHash__hash", passwordHash).
		RelatedSel().
		One(user)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Unknown error for phone: %s", phone))
	}
	_, err = b.ormer.LoadRelated(user, "AdminInfo")
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Unknown error when loading related for phone: %s", phone))
	}
	user.AdminInfo.User = user
	return user.AdminInfo, nil
}

func (b *barCrudImpl) GetById(id uint64) (*AdminInfo, error) {
	profile := &AdminInfo{}
	err := b.ormer.
		QueryTable(profile).
		RelatedSel().
		Filter("id", id).
		One(profile)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err == orm.ErrMissPK {
		return nil, errors.Wrap(err, fmt.Sprintf("Miss primary key for id %d", id))
	}
	return profile, err
}

func (b *barCrudImpl) Update(profile *UpdateAdminInfo) (*AdminInfo, error) {
	params := orm.Params{}
	if profile.Name != nil {
		params["Name"] = profile.Name
	}
	if profile.Surname != nil {
		params["Surname"] = profile.Surname
	}
	if profile.Patronymic != nil {
		params["Patronymic"] = profile.Patronymic
	}
	if profile.Email != nil {
		params["Email"] = profile.Email
	}
	if profile.Phone != nil {
		params["Phone"] = profile.Phone
	}
	if len(params) == 0 {
		return b.GetById(profile.Id)
	}
	num, err := b.ormer.QueryTable(&AdminInfo{}).
		Filter("id", profile.Id).
		Update(params)
	if num == 0 {
		return nil, NoRowsUpdated
	}
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating profile")
	}
	return b.GetById(profile.Id)
}
