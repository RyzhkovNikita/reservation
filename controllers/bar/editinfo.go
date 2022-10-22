package bar

import (
	"barckend/controllers/base"
	"barckend/controllers/requests/bodies"
	"barckend/controllers/requests/input"
	"barckend/crud"
	"barckend/mapping"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
)

type EditController struct {
	base.Controller
}

func (c EditController) EditBar() {
	in := &input.BarIdInPathInput{}
	err := input.ParseInput(c.Ctx.Input, in)
	if err != nil {
		c.BadRequest("Invalid input")
	}
	updateInfo := &bodies.UpdateBar{}
	err = bodies.Require(updateInfo, c.Ctx.Input.RequestBody)
	if err != nil {
		c.BadRequest(err.Error())
	}
	barInfoDb, err := crud.GetBarCrud().GetBarById(uint64(in.BarId))
	if err != nil {
		c.InternalServerError(err)
	}
	if barInfoDb == nil {
		c.NotFound("No bar with provided id: " + strconv.FormatInt(int64(in.BarId), 10))
	}
	ownerId := c.GetUser().OwnerInfo.Id
	if barInfoDb.OwnerInfo.Id != ownerId {
		c.Forbidden()
	}
	if updateInfo.Name != nil {
		isOccupied, er := crud.GetBarCrud().IsNameOccupiedByAnotherOwner(ownerId, *updateInfo.Name)
		if er != nil {
			c.InternalServerError(er)
		}
		if isOccupied {
			c.BadRequest(fmt.Sprintf("Bar name %s is already used by another owner", *updateInfo.Name))
		}
	}
	var workHoursIn []bodies.WorkHours
	if updateInfo.WorkHours != nil {
		workHoursIn = *updateInfo.WorkHours
	}
	updatedBar, err := crud.GetBarCrud().UpdateBar(&crud.UpdateBar{
		Id:          barInfoDb.Id,
		Email:       updateInfo.Email,
		Phone:       updateInfo.Phone,
		Name:        updateInfo.Name,
		Description: updateInfo.Description,
		Address:     updateInfo.Address,
		IsVisible:   updateInfo.IsVisible,
		WorkHours:   mapping.Mapper.WorkHoursListInToDb(barInfoDb.Id, workHoursIn),
	}, barInfoDb)
	if err != nil {
		c.InternalServerError(err)
	}
	if updatedBar == nil {
		c.InternalServerError(errors.New("no bar found after update"))
	}
	c.Data["json"] = mapping.Mapper.BarInfoDbToNet(updatedBar)
	c.ServeJSONInternal()
}
