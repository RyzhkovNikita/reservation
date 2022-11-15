package bar

import (
	"barckend/controllers/base"
	"barckend/controllers/requests/bodies"
	"barckend/crud"
	"barckend/mapping"
	"fmt"
)

type CreateController struct {
	base.Controller
}

func (c *CreateController) CreateBar() {
	barInfo := &bodies.CreateBarInfo{}
	if err := bodies.Require(barInfo, c.Ctx.Input.RequestBody); err != nil {
		c.BadRequest(err.Error())
	}
	isOccupied, err := crud.GetBarCrud().IsNameOccupiedByAnotherOwner(c.GetUser().Id, barInfo.Name)
	if err != nil {
		c.InternalServerError(err)
	}
	if isOccupied {
		c.BadRequest(fmt.Sprintf("Bar name %s is already used by another owner", barInfo.Name))
	}
	createdBar, err := crud.GetBarCrud().InsertBar(&crud.Bar{
		Email:       barInfo.Email,
		Address:     barInfo.Address,
		Name:        barInfo.Name,
		Description: barInfo.Description,
		Phone:       barInfo.Phone,
		OwnerInfo:   &crud.OwnerInfo{Id: c.GetUser().OwnerInfo.Id},
	})
	if err != nil {
		c.InternalServerError(err)
	}
	workHoursListInToDb, err := mapping.Mapper.WorkHoursListInToDb(createdBar.Id, barInfo.WorkHours)
	if err != nil {
		c.InternalServerError(err)
	}
	workHours, err := crud.GetBarCrud().InsertWorkHours(
		createdBar.Id,
		workHoursListInToDb,
	)
	if err != nil {
		c.InternalServerError(err)
	}
	createdBar.WorkHours = workHours
	barInfoResponse, err := mapping.Mapper.BarInfoDbToNet(createdBar)
	if err != nil {
		c.InternalServerError(err)
	}
	c.ServeJSONInternal(barInfoResponse)
}
