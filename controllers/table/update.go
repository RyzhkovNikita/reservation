package table

import (
	"barckend/controllers/base"
	"barckend/controllers/requests/bodies"
	"barckend/controllers/requests/input"
	"barckend/crud"
	"barckend/mapping"
	"strconv"
)

type UpdateController struct {
	base.Controller
}

func (c *UpdateController) UpdateTableById() {
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

	updateInfo := &bodies.UpdateTable{}
	if err = bodies.Require(updateInfo, c.Ctx.Input.RequestBody); err != nil {
		c.BadRequest(err.Error())
	}
	if updateInfo.Name != nil {
		tableByName, err := crud.GetBarCrud().GetTableByName(tableInfo.BarInfo.Id, *updateInfo.Name)
		if err != nil {
			c.InternalServerError(err)
		}
		if tableByName != nil && tableByName.Id != tableInfo.Id {
			c.BadRequest("Table name is already occupied by another table")
		}
	}
	updatedTable, err := crud.GetBarCrud().UpdateTable(
		&crud.UpdateTable{
			Id:          tableId,
			Name:        updateInfo.Name,
			Description: updateInfo.Description,
			Capacity:    updateInfo.Capacity,
		},
		tableInfo,
	)
	if err != nil {
		c.InternalServerError(err)
	}
	c.ServeJSONInternal(mapping.Mapper.TableInfoDbToNet(updatedTable))
}
