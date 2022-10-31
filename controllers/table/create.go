package table

import (
	"barckend/controllers/base"
	"barckend/controllers/requests/bodies"
	"barckend/crud"
	"barckend/mapping"
)

type CreateController struct {
	base.Controller
}

func (c *CreateController) CreateTable() {
	tableInfo := &bodies.CreateTable{}
	if err := bodies.Require(tableInfo, c.Ctx.Input.RequestBody); err != nil {
		c.BadRequest(err.Error())
	}
	inBarId := uint64(tableInfo.BarId)
	c.BarAccessCheck(inBarId)
	tableByName, err := crud.GetBarCrud().GetTableByName(inBarId, tableInfo.Name)
	if err != nil {
		c.InternalServerError(err)
	}
	if tableByName != nil {
		c.BadRequest("Table name is used by another table")
	}
	insertedTable, err := crud.GetBarCrud().InsertTable(&crud.Table{
		Name:        tableInfo.Name,
		Description: *tableInfo.Description,
		Capacity:    uint(tableInfo.Capacity),
		BarInfo:     &crud.Bar{Id: inBarId},
	})
	if err != nil {
		c.InternalServerError(err)
	}
	c.ServeJSONInternal(mapping.Mapper.TableInfoDbToNet(insertedTable))
}
