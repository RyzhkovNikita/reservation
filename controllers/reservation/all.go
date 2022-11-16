package reservation

import (
	"barckend/controllers/base"
	"barckend/controllers/requests/input"
	"barckend/crud"
	"barckend/mapping"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

type AllReservationsController struct {
	base.Controller
}

func (c *AllReservationsController) GetAll() {
	in := &input.TableIdInQueryInput{}
	if err := input.ParseInput(c.Ctx.Input, in); err != nil {
		c.BadRequest("Invalid input")
	}
	tableId := uint64(in.TableId)
	table, err := crud.GetBarCrud().GetTableById(tableId)
	if err != nil {
		c.InternalServerError(err)
	}
	if table == nil {
		c.NotFound("No table with this id")
	}
	inDates := &input.DatesInQueryInput{}
	if err = input.ParseInput(c.Ctx.Input, inDates); err != nil {
		c.BadRequest("Invalid input")
	}
	dates, err := inDates.GetDates()
	if err != nil {
		c.BadRequest("Incorrect dates input")
	}
	barId := table.BarInfo.Id
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
	barInfo, err := crud.GetBarCrud().GetBarById(barId)
	if err != nil {
		c.InternalServerError(err)
	}
	reservations := make([]*crud.Reservation, 0, 10)
	for _, date := range dates {
		reservationsForTableAndDate, err := crud.GetReservCrud().GetReservationsForTableAndDate(tableId, date)
		if err != nil {
			c.InternalServerError(err)
		}
		reservations = append(reservations, reservationsForTableAndDate...)
	}
	table.BarInfo = barInfo
	for _, res := range reservations {
		res.Table = table
	}
	netReservations, err := mapping.Mapper.ReservationListDbToNet(reservations)
	if err != nil {
		c.InternalServerError(err)
	}
	c.ServeJSONInternal(netReservations)
}
