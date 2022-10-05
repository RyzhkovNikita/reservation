package controllers

import (
	"barckend/crud"
)

type ProfileController struct {
	BaseController
	Crud crud.Crud
}

//func (c *ProfileController) GetBarProfile() {
//	c.Data["json"] = c.Mapper.AdminDbToNet(c.GetProfile())
//	err := c.ServeJSON()
//	if err != nil {
//		c.InternalServerError(err)
//	}
//}

func (c *ProfileController) UpdateProfile() {

}

func (c *ProfileController) UploadLogo() {

}
