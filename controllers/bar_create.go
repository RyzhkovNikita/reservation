package controllers

import "barckend/crud"

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
	c.BarCrud.IsNameOccupiedByAnotherOwner(c.GetUser().Id, barInfo.Name)
}
