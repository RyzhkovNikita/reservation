package crud

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/pkg/errors"
	"strconv"
)

func initializeBarCrud(ormer orm.Ormer) {
	instance = &barCrudImpl{ormer: ormer}
}

func GetBarCrud() BarCrud {
	if instance == nil {
		panic("Package is not initialized")
	}
	return instance
}

var instance *barCrudImpl = nil

type BarCrud interface {
	InsertBar(bar *Bar) (*Bar, error)
	InsertWorkHours(barId uint64, workHoursList []WorkHours) ([]*WorkHours, error)
	GetBarById(id uint64) (*Bar, error)
	GetWorkHoursForBar(barId uint64) ([]*WorkHours, error)
	IsNameOccupiedByAnotherOwner(ownerId uint64, name string) (bool, error)
	UpdateBar(updateBar *UpdateBar, barInfo *Bar) (*Bar, error)
}

type barCrudImpl struct {
	ormer orm.Ormer
}

func (b *barCrudImpl) InsertBar(bar *Bar) (*Bar, error) {
	id, err := b.ormer.Insert(bar)
	if err != nil {
		return nil, errors.Wrap(
			err,
			fmt.Sprintf("Insert bar error occured. Bar: %v", bar),
		)
	}
	bar.Id = uint64(id)
	return bar, nil
}

func (b *barCrudImpl) InsertWorkHours(barId uint64, workHoursList []WorkHours) ([]*WorkHours, error) {
	_, err := b.ormer.InsertMulti(len(workHoursList), workHoursList)
	if err != nil {
		return nil, errors.Wrap(err, "Error while inserting work hours, barId="+strconv.FormatUint(barId, 10))
	}
	return b.GetWorkHoursForBar(barId)
}

func (b *barCrudImpl) GetWorkHoursForBar(barId uint64) ([]*WorkHours, error) {
	var workHourListActual = make([]*WorkHours, 7)
	_, err := b.ormer.QueryTable(&WorkHours{}).Filter("bar_id", barId).All(&workHourListActual)
	if err != nil {
		return nil, errors.Wrap(err, "Error while querying work hours, barId="+strconv.FormatUint(barId, 10))
	}
	return workHourListActual, nil
}

func (b *barCrudImpl) GetBarById(id uint64) (*Bar, error) {
	barInfo := &Bar{}
	err := b.ormer.QueryTable(barInfo).
		Filter("id", id).
		One(barInfo)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err == orm.ErrMissPK {
		return nil, errors.Wrap(err, fmt.Sprintf("Miss primary key for id %d", id))
	}
	return barInfo, nil
}

func (b *barCrudImpl) IsNameOccupiedByAnotherOwner(ownerId uint64, name string) (bool, error) {
	count, err := b.ormer.QueryTable(&Bar{}).
		Filter("name", name).
		Exclude("owner_info_id", ownerId).
		Count()
	if err == orm.ErrNoRows {
		return false, nil
	} else if err == orm.ErrMissPK {
		return false, errors.Wrap(err, fmt.Sprintf("Miss primary key"))
	}
	return count > 0, nil
}

func (b *barCrudImpl) UpdateBar(updateBar *UpdateBar, barInfo *Bar) (*Bar, error) {
	if barInfo == nil {
		barInfoDb, err := b.GetBarById(updateBar.Id)
		if err != nil {
			return nil, err
		}
		if barInfoDb == nil {
			return nil, errors.New(fmt.Sprintf("No bar was found with provided id: %d", updateBar.Id))
		}
	} else if updateBar.Id != barInfo.Id {
		return nil, errors.New("Provided bar info and update info has different ids")
	}
	params := orm.Params{}
	addStringIfNeeded(&params, "Email", updateBar.Email, barInfo.Email)
	addStringIfNeeded(&params, "Phone", updateBar.Phone, barInfo.Phone)
	addStringIfNeeded(&params, "Address", updateBar.Address, barInfo.Address)
	addStringIfNeeded(&params, "Name", updateBar.Name, barInfo.Name)
	addStringIfNeeded(&params, "Description", updateBar.Description, barInfo.Description)
	addStringIfNeeded(&params, "LogoUrl", updateBar.LogoUrl, barInfo.LogoUrl)
	addBoolIfNeeded(&params, "IsVisible", updateBar.IsVisible, barInfo.IsVisible)
	if len(params) == 0 {
		return barInfo, nil
	}
	num, err := b.ormer.QueryTable(&Bar{}).
		Filter("id", updateBar.Id).
		Update(params)
	if num == 0 {
		return nil, NoRowsUpdated
	}
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating bar")
	}
	updatedBar, err := b.GetBarById(updateBar.Id)
	if err != nil {
		return nil, errors.Wrap(err, "Error while fetching bar")
	}
	return updatedBar, nil
}
