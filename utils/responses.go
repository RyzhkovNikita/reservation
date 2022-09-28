package utils

import (
	beego "github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) InternalServerError(err error) {
	c.CustomAbort(500, "Internal server error\n"+err.Error())
}

func (c *BaseController) BadRequest(message string) {
	c.CustomAbort(400, message)
}
