package table

import (
	"barckend/controllers/base"
	"barckend/controllers/requests/input"
	"barckend/crud"
	"barckend/mapping"
)

type DeleteController struct {
	base.Controller
}

func (c *DeleteController) DeleteTableById() {
	in := &input.TableIdInPathInput{}
	if err := input.ParseInput(c.Ctx.Input, in); err != nil {
		c.BadRequest("Invalid input")
	}
	tableId := uint64(in.TableId)
	tableInfo, err := crud.GetBarCrud().GetTableById(tableId)
	if err != nil {
		c.InternalServerError(err)
	}
	if tableInfo == nil {
		c.BadRequest("No table with provided id")
	}
	c.BarAccessCheck(tableInfo.BarInfo.Id)
	err = crud.GetBarCrud().RemoveTable(tableId)
	if err != nil {
		c.InternalServerError(err)
	}
	c.ServeJSONInternal(mapping.Mapper.TableInfoDbToNet(tableInfo))
}
