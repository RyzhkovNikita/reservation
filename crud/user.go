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
	force := true
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

type AdminCrud interface {
	GetById(id uint64) (*User, error)
	IsEmailOccupied(email string, ownerId *uint64) (bool, error)
	IsPhoneOccupied(phone string, ownerId *uint64) (bool, error)
	InsertAdmin(profile *AdminInfo, passwordHash string) (*AdminInfo, error)
	InsertOwner(profile *OwnerInfo, passwordHash string) (*OwnerInfo, error)
	UpdateAdmin(profile *UpdateAdminInfo) (*AdminInfo, error)
	UpdateOwner(profile *UpdateOwnerInfo) (*OwnerInfo, error)
	CheckCredentialsEmail(email string, passwordHash string) (*User, error)
	CheckCredentialsPhone(phone string, passwordHash string) (*User, error)
}

type userCrudImpl struct {
	ormer orm.Ormer
}

func (b *userCrudImpl) IsEmailOccupied(email string, ownerId *uint64) (bool, error) {
	user := &User{}
	exprAdmin := b.ormer.QueryTable(user).
		RelatedSel().
		Filter("AdminInfo__email", email)
	if ownerId != nil {
		exprAdmin = exprAdmin.Exclude("id", ownerId)
	}
	exprOwner := b.ormer.QueryTable(user).
		RelatedSel().
		Filter("OwnerInfo__email", email)
	if ownerId != nil {
		exprAdmin = exprAdmin.Exclude("id", ownerId)
	}
	return exprAdmin.Exist() || exprOwner.Exist(), nil
}

func (b *userCrudImpl) IsPhoneOccupied(phone string, ownerId *uint64) (bool, error) {
	user := &User{}
	exprAdmin := b.ormer.QueryTable(user).
		RelatedSel().
		Filter("AdminInfo__phone", phone)
	if ownerId != nil {
		exprAdmin = exprAdmin.Exclude("id", ownerId)
	}
	exprOwner := b.ormer.QueryTable(user).
		RelatedSel().
		Filter("OwnerInfo__phone", phone)
	if ownerId != nil {
		exprAdmin = exprAdmin.Exclude("id", ownerId)
	}
	return exprAdmin.Exist() || exprOwner.Exist(), nil
}

func (b *userCrudImpl) InsertAdmin(profile *AdminInfo, passwordHash string) (*AdminInfo, error) {
	user := &User{
		Role:     Admin,
		IsActive: true,
	}
	passwordHashModel := &PasswordHash{
		Hash: passwordHash,
		User: user,
	}
	passwordHashModel.User = user
	_, err := b.ormer.Insert(user)
	if err != nil {
		return nil, errors.Wrap(
			err,
			fmt.Sprintf("Insert user error occured. User: %v", user),
		)
	}
	_, err = b.ormer.Insert(passwordHashModel)
	if err != nil {
		return nil, errors.Wrap(
			err,
			fmt.Sprintf("Insert password error occured. Password hash: %v", passwordHashModel),
		)
	}
	profile.User = user
	_, err = b.ormer.Insert(profile)
	if err != nil {
		return nil, errors.Wrap(
			err,
			fmt.Sprintf("Insert admin info error occured. AdminInfo: %v", profile),
		)
	}
	return profile, nil
}

func (b *userCrudImpl) InsertOwner(profile *OwnerInfo, passwordHash string) (*OwnerInfo, error) {
	user := &User{
		Role:     Owner,
		IsActive: true,
	}
	passwordHashModel := &PasswordHash{
		Hash: passwordHash,
		User: user,
	}
	passwordHashModel.User = user
	_, err := b.ormer.Insert(user)
	if err != nil {
		return nil, errors.Wrap(
			err,
			fmt.Sprintf("Insert user error occured. User: %v", user),
		)
	}
	_, err = b.ormer.Insert(passwordHashModel)
	if err != nil {
		return nil, errors.Wrap(
			err,
			fmt.Sprintf("Insert password error occured. Password hash: %v", passwordHashModel),
		)
	}
	profile.User = user
	_, err = b.ormer.Insert(profile)
	if err != nil {
		return nil, errors.Wrap(
			err,
			fmt.Sprintf("Insert admin info error occured. AdminInfo: %v", profile),
		)
	}
	return profile, nil
}

func (b *userCrudImpl) CheckCredentialsEmail(email string, passwordHash string) (*User, error) {
	user := &User{}
	err := b.ormer.
		QueryTable(user).
		Filter("AdminInfo__email", email).
		Filter("PasswordHash__hash", passwordHash).
		RelatedSel().
		One(user)
	if err == nil {
		_, err = b.ormer.LoadRelated(user, "AdminInfo")
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Unknown error when loading related for email: %s", email))
		}
		user.AdminInfo.User = user
		return user, nil
	}
	err = b.ormer.
		QueryTable(user).
		Filter("OwnerInfo__email", email).
		Filter("PasswordHash__hash", passwordHash).
		RelatedSel().
		One(user)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Unknown error for email: %s", email))
	}
	_, err = b.ormer.LoadRelated(user, "OwnerInfo")
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Unknown error when loading related for email: %s", email))
	}
	user.OwnerInfo.User = user
	return user, nil
}

func (b *userCrudImpl) CheckCredentialsPhone(phone string, passwordHash string) (*User, error) {
	user := &User{}
	err := b.ormer.
		QueryTable(user).
		Filter("AdminInfo__phone", phone).
		Filter("PasswordHash__hash", passwordHash).
		RelatedSel().
		One(user)
	if err == nil {
		_, err = b.ormer.LoadRelated(user, "AdminInfo")
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Unknown error when loading related for phone: %s", phone))
		}
		user.AdminInfo.User = user
		return user, nil
	}
	err = b.ormer.
		QueryTable(user).
		Filter("OwnerInfo__email", phone).
		Filter("PasswordHash__hash", passwordHash).
		RelatedSel().
		One(user)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Unknown error for phone: %s", phone))
	}
	_, err = b.ormer.LoadRelated(user, "OwnerInfo")
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Unknown error when loading related for phone: %s", phone))
	}
	user.OwnerInfo.User = user
	return user, nil
}

func (b *userCrudImpl) GetById(id uint64) (*User, error) {
	user := &User{}
	err := b.ormer.
		QueryTable(user).
		Filter("id", id).
		One(user)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err == orm.ErrMissPK {
		return nil, errors.Wrap(err, fmt.Sprintf("Miss primary key for id %d", id))
	}
	if user.IsAdmin() {
		_, err = b.ormer.LoadRelated(user, "AdminInfo")
	} else if user.IsOwner() {
		_, err = b.ormer.LoadRelated(user, "OwnerInfo")
	}
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Unknown error when loading related for id: %d", id))
	}
	return user, err
}

func (b *userCrudImpl) UpdateAdmin(profile *UpdateAdminInfo) (*AdminInfo, error) {
	user, err := b.GetById(profile.Id)
	if err != nil {
		return nil, errors.Wrap(err, "Error while fetching profile")
	}
	if user.Role != Admin {
		return nil, errors.New("User has wrong role")
	}
	params := orm.Params{}
	addStringIfNeeded(&params, "Name", profile.Name, user.AdminInfo.Name)
	addStringIfNeeded(&params, "Surname", profile.Surname, user.AdminInfo.Surname)
	addStringIfNeeded(&params, "Patronymic", profile.Patronymic, user.AdminInfo.Patronymic)
	addStringIfNeeded(&params, "Email", profile.Email, user.AdminInfo.Email)
	addStringIfNeeded(&params, "Phone", profile.Phone, user.AdminInfo.Phone)
	if len(params) == 0 {
		return user.AdminInfo, nil
	}
	num, err := b.ormer.QueryTable(&AdminInfo{}).
		Filter("user_id", profile.Id).
		Update(params)
	if num == 0 {
		return nil, NoRowsUpdated
	}
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating profile")
	}
	user, err = b.GetById(profile.Id)
	if err != nil {
		return nil, errors.Wrap(err, "Error while fetching profile")
	}
	return user.AdminInfo, nil
}

func addStringIfNeeded(params *orm.Params, key string, newValue *string, oldValue string) {
	if newValue == nil {
		return
	}
	if *newValue == oldValue {
		return
	}
	(*params)[key] = *newValue
}

func addBoolIfNeeded(params *orm.Params, key string, newValue *bool, oldValue bool) {
	if newValue == nil {
		return
	}
	if *newValue == oldValue {
		return
	}
	(*params)[key] = newValue
}

func (b *userCrudImpl) UpdateOwner(profile *UpdateOwnerInfo) (*OwnerInfo, error) {
	user, err := b.GetById(profile.Id)
	if err != nil {
		return nil, errors.Wrap(err, "Error while fetching profile")
	}
	if user.Role != Owner {
		return nil, errors.New("User has wrong role")
	}
	params := orm.Params{}
	addStringIfNeeded(&params, "Name", profile.Name, user.OwnerInfo.Name)
	addStringIfNeeded(&params, "Surname", profile.Surname, user.OwnerInfo.Surname)
	addStringIfNeeded(&params, "Patronymic", profile.Patronymic, user.OwnerInfo.Patronymic)
	addStringIfNeeded(&params, "Email", profile.Email, user.OwnerInfo.Email)
	addStringIfNeeded(&params, "Phone", profile.Phone, user.OwnerInfo.Phone)
	if len(params) == 0 {
		return user.OwnerInfo, nil
	}
	num, err := b.ormer.QueryTable(&OwnerInfo{}).
		Filter("user_id", profile.Id).
		Update(params)
	if num == 0 {
		return nil, NoRowsUpdated
	}
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating profile")
	}
	user, err = b.GetById(profile.Id)
	if err != nil {
		return nil, errors.Wrap(err, "Error while fetching profile")
	}
	return user.OwnerInfo, nil
}
