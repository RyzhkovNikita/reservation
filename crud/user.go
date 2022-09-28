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
	orm.RegisterModel(new(Profile), new(Credentials))
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
	GetById(id int) (*Profile, error)
	Insert(profile *Profile) (*Profile, error)
	GetByEmail(email string) (*Profile, error)
	CheckCredentials(email string, passwordHash string) (*Profile, error)
}

type barCrudImpl struct {
	ormer orm.Ormer
}

func (crud *barCrudImpl) GetById(id int) (*Profile, error) {
	profile := &Profile{}
	err := crud.ormer.
		QueryTable(profile).
		RelatedSel().
		Filter("Profile__id", id).
		One(profile)
	if err != nil {
		return nil, err
	}
	return profile, err
}

func (crud *barCrudImpl) Insert(profile *Profile) (*Profile, error) {
	_, err := crud.ormer.Insert(profile.Credentials)
	if err != nil {
		return nil, errors.Wrap(
			err,
			fmt.Sprintf("Insert creds error occured. Profile: %v", profile),
		)
	}
	_, err = crud.ormer.Insert(profile)
	if err != nil {
		return nil, errors.Wrap(
			err,
			fmt.Sprintf("Insert profile error occured. Profile: %v", profile),
		)
	}
	return profile, nil
}

func (crud *barCrudImpl) GetByEmail(email string) (*Profile, error) {
	profile := &Profile{}
	err := crud.ormer.
		QueryTable(profile).
		RelatedSel().
		Filter("Credentials__email", email).
		One(profile)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err == orm.ErrMissPK {
		return nil, errors.Wrap(err, fmt.Sprintf("Miss primary key for email %s", email))
	}
	return profile, nil
}

func (crud *barCrudImpl) CheckCredentials(email string, passwordHash string) (*Profile, error) {
	profile := &Profile{}
	err := crud.ormer.
		QueryTable(profile).
		RelatedSel().
		Filter("Credentials__email", email).
		Filter("Credentials__password_hash", passwordHash).
		One(profile)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err == orm.ErrMissPK {
		return nil, errors.Wrap(err, fmt.Sprintf("Miss primary key for email %s", email))
	}
	return profile, nil
}
