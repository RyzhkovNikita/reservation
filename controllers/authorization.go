package controllers

import (
	crud "barckend/crud"
	beego "github.com/beego/beego/v2/server/web"
)

type RegistrationController struct {
	beego.Controller
	crud crud.BarCrud
}

func (c *RegistrationController) Post() {
	profile, err := c.crud.Insert(&crud.BarProfile{})
	profileId := profile.Id
}
