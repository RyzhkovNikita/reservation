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
	GetBarIdsForOwner(ownerId uint64) ([]uint64, error)
	GetBarForAdmin(adminId uint64) (*Bar, error)
	GetTableByName(barId uint64, tableName string) (*Table, error)
	InsertTable(table *Table) (*Table, error)
	UpdateTable(table *UpdateTable) (*Table, error)
	RemoveTable(tableId uint64) error
	GetAllTables(barId uint64) ([]*Table, error)
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
	var workHourListActual = make([]*WorkHours, 0, 7)
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
	_, err = b.ormer.LoadRelated(barInfo, "WorkHours")
	if err != nil {
		return nil, errors.Wrap(err, "Error while loading related work hours")
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
		barInfo = barInfoDb
	} else if updateBar.Id != barInfo.Id {
		return nil, errors.New("Provided bar info and update info has different ids")
	}
	if len(updateBar.WorkHours) != 0 {
		err := b.UpdateWorkHours(updateBar.Id, updateBar.WorkHours)
		if err != nil {
			return nil, errors.Wrap(err, "Error while updating work hours")
		}
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

func (b *barCrudImpl) UpdateWorkHours(barId uint64, newWorkHours []WorkHours) error {
	ormTransaction, err := b.ormer.Begin()
	if err != nil {
		ormTransaction.Rollback()
		return errors.Wrap(err, "Error when start transaction")
	}
	_, err = ormTransaction.QueryTable(&WorkHours{}).Filter("bar_id", barId).Delete()
	if err != nil {
		ormTransaction.Rollback()
		return errors.Wrap(err, "Error when deleting work hours")
	}
	_, err = ormTransaction.InsertMulti(len(newWorkHours), newWorkHours)
	if err != nil {
		ormTransaction.Rollback()
		return errors.Wrap(err, "Error when inserting new work hours")
	}
	err = ormTransaction.Commit()
	if err != nil {
		ormTransaction.Rollback()
		return errors.Wrap(err, "Error when committing transaction")
	}
	return nil
}

func (b *barCrudImpl) GetBarIdsForOwner(ownerId uint64) ([]uint64, error) {
	var barList = make([]*Bar, 0, 5)
	_, err := b.ormer.QueryTable(&Bar{}).
		Filter("owner_info_id", ownerId).
		All(&barList)
	if err != nil {
		return nil, errors.Wrap(err, "Error while querying bar list, ownerId="+strconv.FormatUint(ownerId, 10))
	}
	ids := make([]uint64, 0, len(barList))
	for _, bar := range barList {
		ids = append(ids, bar.Id)
	}
	return ids, nil
}

func (b *barCrudImpl) GetBarForAdmin(adminId uint64) (*Bar, error) {
	adminInfo := &AdminInfo{}
	err := b.ormer.QueryTable(adminInfo).Filter("id", adminId).One(adminInfo)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err == orm.ErrMissPK {
		return nil, errors.Wrap(err, fmt.Sprintf("Miss primary key for id %d", adminId))
	}
	_, err = b.ormer.LoadRelated(adminInfo, "Bar")
	if err != nil {
		return nil, nil
	}
	return adminInfo.Bar, nil
}

func (b *barCrudImpl) GetTableByName(barId uint64, tableName string) (*Table, error) {
	tableInfo := &Table{}
	err := b.ormer.QueryTable(tableInfo).
		Filter("name", tableName).
		Filter("bar_info_id", barId).
		One(tableInfo)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	return tableInfo, nil
}

func (b *barCrudImpl) InsertTable(table *Table) (*Table, error) {
	id, err := b.ormer.Insert(table)
	if err != nil {
		return nil, errors.Wrap(
			err,
			fmt.Sprintf("Insert table error occured. Table: %v", table),
		)
	}
	table.Id = uint64(id)
	return table, nil
}

func (b *barCrudImpl) UpdateTable(table *UpdateTable) (*Table, error) {
	//TODO implement me
	panic("implement me")
}

func (b *barCrudImpl) RemoveTable(tableId uint64) error {
	//TODO implement me
	panic("implement me")
}

func (b *barCrudImpl) GetAllTables(barId uint64) ([]*Table, error) {
	tableList := make([]*Table, 0, 15)
	_, err := b.ormer.QueryTable(&Table{}).
		Filter("bar_info_id", barId).
		All(&tableList)
	if err != nil {
		return nil, errors.Wrap(err, "Error while querying table list, barId="+strconv.FormatUint(barId, 10))
	}
	return tableList, nil
}
