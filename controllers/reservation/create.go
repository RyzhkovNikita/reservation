package reservation

import (
	"barckend/controllers/base"
	"barckend/controllers/requests/bodies"
	"barckend/crud"
	"barckend/mapping"
	"barckend/timing"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"golang.org/x/exp/slices"
	"time"
)

type CreateController struct {
	base.Controller
}

func (c *CreateController) CreateReservation() {
	reservationIn := &bodies.CreateReservation{}
	if err := bodies.Require(reservationIn, c.Ctx.Input.RequestBody); err != nil {
		c.BadRequest(err.Error())
	}
	formattedFromTime := fmt.Sprintf("%s %s", reservationIn.Date, reservationIn.FromTime)
	reservFromTime, err := timing.GetConverter().StringToDateTime(formattedFromTime)
	if err != nil {
		c.BadRequest("Wrong date or timing format")
	}
	minutes := reservFromTime.Minute()
	if !slices.Contains([]int{0, 15, 30, 45}, minutes) {
		c.BadRequest("Wrong date or timing format")
	}
	if reservFromTime.Before(time.Now().Add(-3 * time.Hour)) {
		c.BadRequest("No availability to create reservation in past")
	}
	var reservToTime time.Time
	formattedToTime := fmt.Sprintf("%s %s", reservationIn.Date, reservationIn.ToTime)
	reservToTime, err = timing.GetConverter().StringToDateTime(formattedToTime)
	if err != nil {
		c.BadRequest("Wrong date or timing format")
	}
	if reservToTime.Before(reservFromTime) {
		reservToTime = reservToTime.Add(24 * time.Hour)
	}
	if reservFromTime.Equal(reservToTime) {
		c.BadRequest("Time to reserv is equal to 0")
	}
	minutes = reservToTime.Minute()
	if !slices.Contains([]int{0, 15, 30, 45}, minutes) {
		c.BadRequest("Wrong date or timing format")
	}
	table, err := crud.GetBarCrud().GetTableById(uint64(reservationIn.TableId))
	if err != nil {
		c.InternalServerError(err)
	}
	if table == nil {
		c.NotFound(fmt.Sprintf("No table with provided id %d", reservationIn.TableId))
	}
	if table.Capacity < uint(reservationIn.PersonCount) {
		c.BadRequest(fmt.Sprintf("Table capacity (%d) is less than person count (%d)", table.Capacity, reservationIn.PersonCount))
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
	var weekdayWorkHours *crud.WorkHours
	weekday := reservFromTime.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	for _, wh := range barInfo.WorkHours {
		if wh.Weekday == uint(weekday) {
			weekdayWorkHours = wh
			break
		}
	}
	if weekdayWorkHours == nil {
		c.InternalServerError(errors.New("no work hours for bar"))
	}
	workHoursFrom := weekdayWorkHours.From.Value().AddDate(reservFromTime.Year(), int(reservFromTime.Month())-1, reservFromTime.Day()-1)
	workHoursTo := weekdayWorkHours.To.Value().AddDate(reservFromTime.Year(), int(reservFromTime.Month())-1, reservFromTime.Day()-1)
	if workHoursTo.Before(workHoursFrom) {
		panic(123)
	}
	afterFromHours := reservFromTime.After(workHoursFrom) || reservFromTime.Equal(workHoursFrom)
	beforeToHours := reservFromTime.Before(workHoursTo)
	if !afterFromHours || !beforeToHours {
		c.BadRequest("Invalid reservation timing - start of reservation is outside of working hours")
	}
	if reservToTime.After(workHoursTo) {
		reservToTime = workHoursTo
	}
	fromTimeField := orm.DateTimeField(reservFromTime)
	toTimeField := orm.DateTimeField(reservToTime)
	dateField := orm.DateField(reservFromTime)
	reservation := crud.Reservation{
		From:        fromTimeField,
		To:          toTimeField,
		Date:        dateField,
		Table:       &crud.Table{Id: table.Id},
		PersonCount: uint8(reservationIn.PersonCount),
	}
	reservations, err := crud.GetReservCrud().GetReservationsForTableAndDate(table.Id, reservFromTime)
	if err != nil {
		c.InternalServerError(err)
	}
	if isInterceptingWith(reservations, reservation) {
		c.BadRequest("Reservation is intercepting with existing")
	}
	insertedReservation, err := crud.GetReservCrud().InsertReservation(&reservation)
	if err != nil {
		c.InternalServerError(err)
	}
	table.BarInfo = barInfo
	insertedReservation.Table = table
	reservationResponse, err := mapping.Mapper.ReservationDbToNet(insertedReservation)
	if err != nil {
		c.InternalServerError(err)
	}
	c.ServeJSONInternal(reservationResponse)
}

func isInterceptingWith(reservations []*crud.Reservation, target crud.Reservation) bool {
	for _, res := range reservations {
		if target.IsInterceptingWith(*res) {
			return true
		}
	}
	return false
}
