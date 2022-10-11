package controllers

import (
	"barckend/crud"
	"fmt"
)

type ProfileController struct {
	BaseController
}

func (c *ProfileController) GetMe() {
	c.Data["json"] = c.Mapper.AdminDbToNet(c.GetProfile())
	err := c.ServeJSON()
	if err != nil {
		c.InternalServerError(err)
	}
}

func (c *ProfileController) PatchMe() {
	admin := &UpdateProfile{}
	if err := Require(admin, c.Ctx.Input.RequestBody); err != nil {
		c.BadRequest(err.Error())
	}
	if admin.Email != nil {
		isOccupied, err := c.Crud.IsEmailOccupied(*admin.Email, &c.GetProfile().Id)
		if err != nil {
			c.InternalServerError(err)
		}
		if isOccupied {
			c.BadRequest(
				fmt.Sprintf("Profile with this email {%s} is already existing", *admin.Email),
			)
		}
	}
	if admin.Phone != nil {
		isOccupied, err := c.Crud.IsPhoneOccupied(*admin.Phone, &c.GetProfile().Id)
		if err != nil {
			c.InternalServerError(err)
		}
		if isOccupied {
			c.BadRequest(
				fmt.Sprintf("Profile with this Phone {%s} is already existing", *admin.Phone),
			)
		}
	}
	adminInfo, err := c.Crud.Update(&crud.UpdateAdminInfo{
		Id:         c.GetProfile().Id,
		Surname:    admin.Surname,
		Name:       admin.Name,
		Patronymic: admin.Patronymic,
		Email:      admin.Email,
		Phone:      admin.Phone,
	})
	if err != nil {
		c.InternalServerError(err)
	}
	if adminInfo == nil {
		c.InternalServerError(fmt.Errorf("no user found after update"))
	}
	c.Data["json"] = c.Mapper.AdminDbToNet(adminInfo)
	err = c.ServeJSON()
	if err != nil {
		c.InternalServerError(err)
	}
}
