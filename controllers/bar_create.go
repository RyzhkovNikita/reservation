package controllers

import (
	"barckend/crud"
	"fmt"
)

type BarCreateController struct {
	BaseController
	BarCrud crud.BarCrud
	Mapper  ModelMapper
}

func (c *BarCreateController) CreateBar() {
	barInfo := &CreateBarInfo{}
	if err := Require(barInfo, c.Ctx.Input.RequestBody); err != nil {
		c.BadRequest(err.Error())
	}
	isOccupied, err := c.BarCrud.IsNameOccupiedByAnotherOwner(c.GetUser().Id, barInfo.Name)
	if err != nil {
		c.InternalServerError(err)
	}
	if isOccupied {
		c.BadRequest(fmt.Sprintf("Bar name %s is already used by another owner", barInfo.Name))
	}
	createdBar, err := c.BarCrud.InsertBar(&crud.Bar{
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
	_, err = c.BarCrud.InsertWorkHours(
		createdBar.Id,
		c.Mapper.WorkHoursListInToDb(createdBar.Id, barInfo.WorkHours),
	)
	c.Data["json"] = c.Mapper.BarInfoDbToNet(createdBar)
	c.ServeJSONInternal()
}
