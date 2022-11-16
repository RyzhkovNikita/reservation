package reservation

import (
	"barckend/controllers/base"
	"barckend/controllers/requests/input"
	"barckend/crud"
	"barckend/mapping"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

type DeleteController struct {
	base.Controller
}

func (c *DeleteController) DeleteById() {
	in := &input.ReservationIdInPathInput{}
	if err := input.ParseInput(c.Ctx.Input, in); err != nil {
		c.BadRequest("Invalid input")
	}
	reservationId := uint64(in.ReservationId)
	reservation, err := crud.GetReservCrud().GetReservationById(reservationId)
	if err != nil {
		c.InternalServerError(err)
	}
	if reservation == nil {
		c.BadRequest("Reservation not found")
	}
	barId := reservation.Table.BarInfo.Id
	user := c.GetUser()
	if user.IsOwner() {
		barIds, err := crud.GetBarCrud().GetBarIdsForOwner(user.OwnerInfo.Id)
		if err != nil {
			c.InternalServerError(err)
		}
		if !slices.Contains(barIds, barId) {
			c.Forbidden()
		}
	} else if user.IsAdmin() {
		barForAdmin, err := crud.GetBarCrud().GetBarForAdmin(user.AdminInfo.Id)
		if err != nil {
			c.InternalServerError(err)
		}
		if barForAdmin.Id != barId {
			c.Forbidden()
		}
	} else {
		c.InternalServerError(errors.New(""))
	}
	err = crud.GetReservCrud().DeleteReservationById(reservationId)
	if err != nil {
		c.InternalServerError(err)
	}
	reservationDbToNet, err := mapping.Mapper.ReservationDbToNet(reservation)
	if err != nil {
		c.InternalServerError(err)
	}
	c.ServeJSONInternal(reservationDbToNet)
}
