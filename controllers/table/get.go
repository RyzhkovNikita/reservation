package table

import (
	"barckend/controllers/base"
	"barckend/controllers/requests/input"
	"barckend/crud"
	"barckend/mapping"
	"strconv"
)

type GetController struct {
	base.Controller
}

func (c *GetController) GetTableById() {
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
		c.NotFound("No table with provided id=" + strconv.FormatUint(tableId, 10))
	}
	c.BarAccessCheck(tableInfo.BarInfo.Id)
	c.ServeJSONInternal(mapping.Mapper.TableInfoDbToNet(tableInfo))
}
