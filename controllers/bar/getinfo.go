package bar

import (
	"barckend/controllers/base"
	"barckend/controllers/requests/input"
	"barckend/crud"
	"barckend/mapping"
	"strconv"
)

type InfoController struct {
	base.Controller
}

func (c *InfoController) GetBarInformation() {
	in := &input.BarIdInPathInput{}
	err := input.ParseInput(c.Ctx.Input, in)
	if err != nil {
		c.BadRequest("Invalid input")
	}
	barInfoDb, err := crud.GetBarCrud().GetBarById(uint64(in.BarId))
	if err != nil {
		c.InternalServerError(err)
	}
	if barInfoDb == nil {
		c.NotFound("No bar with provided id: " + strconv.FormatInt(int64(in.BarId), 10))
	}
	barInfoResponse, err := mapping.Mapper.BarInfoDbToNet(barInfoDb)
	if err != nil {
		c.InternalServerError(err)
	}
	c.ServeJSONInternal(barInfoResponse)
}
