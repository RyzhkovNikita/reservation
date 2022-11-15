package crud

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

func initReservCrud(ormer orm.Ormer) {
	reservationCrudInstance = &reservationCrudImpl{ormer: ormer}
}

func GetReservCrud() ReservationCrud {
	if reservationCrudInstance == nil {
		panic("Package is not initialized")
	}
	return *reservationCrudInstance
}

type ReservationCrud interface {
	GetReservationsForTableAndDate(tableId uint64, timee time.Time) ([]*Reservation, error)
	InsertReservation(reservation *Reservation) (*Reservation, error)
}

var reservationCrudInstance *reservationCrudImpl

type reservationCrudImpl struct {
	ormer orm.Ormer
}

func (r reservationCrudImpl) GetReservationsForTableAndDate(
	tableId uint64,
	timee time.Time,
) ([]*Reservation, error) {
	reservations := make([]*Reservation, 0, 5)
	_, err := r.ormer.QueryTable("reservation").
		Filter("table_id", tableId).
		Filter("date", timee).
		All(&reservations)
	if err != nil {
		return nil, errors.Wrap(err, "Error while querying reservations, table_id="+strconv.FormatUint(tableId, 10))
	}
	return reservations, nil
}

func (r reservationCrudImpl) InsertReservation(reservation *Reservation) (*Reservation, error) {
	id, err := r.ormer.Insert(reservation)
	if err != nil {
		return nil, errors.Wrap(
			err,
			fmt.Sprintf("Insert reservations error occured. Reservation: %v", *reservation),
		)
	}
	reservation.Id = uint64(id)
	return reservation, nil
}
