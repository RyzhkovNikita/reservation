package crud

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/pkg/errors"
)

var BarDb BarCrud

type BarCrud interface {
	Create(bar *Bar) (*Bar, error)
	GetById(id uint64) (*Bar, error)
	IsNameOccupiedByAnotherOwner(ownerId uint64, name string) (bool, error)
	Update(updateBar UpdateBar) (*Bar, error)
}

type barCrudImpl struct {
	ormer orm.Ormer
}

func (b *barCrudImpl) Create(bar *Bar) (*Bar, error) {
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

func (b *barCrudImpl) GetById(id uint64) (*Bar, error) {
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
		Filter("owner_info_id", ownerId).
		Count()
	if err == orm.ErrNoRows {
		return false, nil
	} else if err == orm.ErrMissPK {
		return false, errors.Wrap(err, fmt.Sprintf("Miss primary key"))
	}
	return count > 0, nil
}

func (b *barCrudImpl) Update(updateBar UpdateBar) (*Bar, error) {
	//TODO implement me
	panic("implement me")
}
