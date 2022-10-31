package table

import (
	"barckend/controllers/base"
	"barckend/controllers/requests/input"
	"barckend/crud"
	"barckend/mapping"
)

type GetAllController struct {
	base.Controller
}

func (c *GetAllController) GetAllTables() {
	in := &input.BarIdInQueryInput{}
	if err := input.ParseInput(c.Ctx.Input, in); err != nil {
		c.BadRequest("Invalid input")
	}
	inBarId := uint64(in.BarId)
	c.BarAccessCheck(inBarId)
	tables, err := crud.GetBarCrud().GetAllTables(inBarId)
	if err != nil {
		c.InternalServerError(err)
	}
	c.ServeJSONInternal(mapping.Mapper.TableInfoListDbToNet(tables))
}
