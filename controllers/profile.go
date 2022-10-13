package controllers

import (
	"barckend/crud"
	"fmt"
	"github.com/pkg/errors"
)

type ProfileController struct {
	BaseController
	Mapper ModelMapper
}

func (c *ProfileController) GetMe() {
	user := c.GetUser()
	if user.IsAdmin() {
		c.Data["json"] = c.Mapper.AdminDbToNet(user.AdminInfo)
	} else if user.IsOwner() {
		c.Data["json"] = c.Mapper.OwnerDbToNet(user.OwnerInfo)
	}
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
	userMe := c.GetUser()
	if admin.Email != nil {
		isOccupied, err := c.Crud.IsEmailOccupied(*admin.Email, &userMe.Id)
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
		isOccupied, err := c.Crud.IsPhoneOccupied(*admin.Phone, &userMe.Id)
		if err != nil {
			c.InternalServerError(err)
		}
		if isOccupied {
			c.BadRequest(
				fmt.Sprintf("Profile with this Phone {%s} is already existing", *admin.Phone),
			)
		}
	}
	var err error
	var adminInfo *crud.AdminInfo
	var ownerInfo *crud.OwnerInfo
	if userMe.IsAdmin() {
		adminInfo, err = c.Crud.UpdateAdmin(&crud.UpdateAdminInfo{
			Id:         userMe.Id,
			Surname:    admin.Surname,
			Name:       admin.Name,
			Patronymic: admin.Patronymic,
			Email:      admin.Email,
			Phone:      admin.Phone,
		})
	} else {
		ownerInfo, err = c.Crud.UpdateOwner(&crud.UpdateOwnerInfo{
			Id:         userMe.Id,
			Surname:    admin.Surname,
			Name:       admin.Name,
			Patronymic: admin.Patronymic,
			Email:      admin.Email,
			Phone:      admin.Phone,
		})
	}
	if err != nil {
		c.InternalServerError(err)
	}
	if adminInfo == nil && ownerInfo == nil {
		c.InternalServerError(errors.New("no user found after update"))
	}
	if userMe.IsAdmin() {
		c.Data["json"] = c.Mapper.AdminDbToNet(adminInfo)
	} else if userMe.IsOwner() {
		c.Data["json"] = c.Mapper.OwnerDbToNet(ownerInfo)
	} else {
		c.InternalServerError(errors.New("WTF"))
	}
	err = c.ServeJSON()
	if err != nil {
		c.InternalServerError(err)
	}
}
