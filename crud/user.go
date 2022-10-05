package crud

import (
	"barckend/conf"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/pkg/errors"

	//_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var Db Crud

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

type Crud interface {
	IsEmailExist(email string) (bool, error)
	IsPhoneExist(phone string) (bool, error)
	Insert(profile *AdminInfo) (*AdminInfo, error)
	CheckCredentialsEmail(email string, passwordHash string) (*AdminInfo, error)
	CheckCredentialsPhone(phone string, passwordHash string) (*AdminInfo, error)
}

type barCrudImpl struct {
	ormer orm.Ormer
}

func (b *barCrudImpl) IsEmailExist(email string) (bool, error) {
	adminInfo := &AdminInfo{}
	err := b.ormer.QueryTable(adminInfo).
		Filter("Email", email).
		One(adminInfo)
	if err == orm.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, errors.Wrap(err, fmt.Sprintf("Unknown error"))
	}
	return true, nil
}

func (b *barCrudImpl) IsPhoneExist(phone string) (bool, error) {
	adminInfo := &AdminInfo{}
	err := b.ormer.QueryTable(adminInfo).
		Filter("Phone", phone).
		One(adminInfo)
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

//func (crud *barCrudImpl) GetById(id int64) (*Profile, error) {
//	profile := &Profile{}
//	err := crud.ormer.
//		QueryTable(profile).
//		RelatedSel().
//		Filter("Profile__id", id).
//		One(profile)
//	if err == orm.ErrNoRows {
//		return nil, nil
//	} else if err == orm.ErrMissPK {
//		return nil, errors.Wrap(err, fmt.Sprintf("Miss primary key for id %d", id))
//	}
//	return profile, err
//}
//
//func (crud *barCrudImpl) Insert(profile *Profile) (*Profile, error) {
//	_, err := crud.ormer.Insert(profile.Credentials)
//	if err != nil {
//		return nil, errors.Wrap(
//			err,
//			fmt.Sprintf("Insert creds error occured. Profile: %v", profile),
//		)
//	}
//	_, err = crud.ormer.Insert(profile)
//	if err != nil {
//		return nil, errors.Wrap(
//			err,
//			fmt.Sprintf("Insert profile error occured. Profile: %v", profile),
//		)
//	}
//	return profile, nil
//}

//func (crud *barCrudImpl) GetByEmail(email string) (*Profile, error) {
//	profile := &Profile{}
//	err := crud.ormer.
//		QueryTable(profile).
//		RelatedSel().
//		Filter("Credentials__email", email).
//		One(profile)
//	if err == orm.ErrNoRows {
//		return nil, nil
//	} else if err == orm.ErrMissPK {
//		return nil, errors.Wrap(err, fmt.Sprintf("Miss primary key for email %s", email))
//	}
//	return profile, nil
//}
//
//func (crud *barCrudImpl) CheckCredentials(email string, passwordHash string) (*Profile, error) {
//	profile := &Profile{}
//	err := crud.ormer.
//		QueryTable(profile).
//		RelatedSel().
//		Filter("Credentials__email", email).
//		Filter("Credentials__password_hash", passwordHash).
//		One(profile)
//	if err == orm.ErrNoRows {
//		return nil, nil
//	} else if err == orm.ErrMissPK {
//		return nil, errors.Wrap(err, fmt.Sprintf("Miss primary key for email %s", email))
//	}
//	return profile, nil
//}
//
