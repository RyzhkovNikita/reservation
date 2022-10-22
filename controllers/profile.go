package controllers

import (
	"barckend/controllers/base"
	"barckend/controllers/requests/bodies"
	"barckend/crud"
	"barckend/mapping"
	"fmt"
	"github.com/pkg/errors"
)

type ProfileController struct {
	base.Controller
}

func (c *ProfileController) GetMe() {
	user := c.GetUser()
	if user.IsAdmin() {
		c.Data["json"] = mapping.Mapper.AdminDbToNet(user.AdminInfo)
	} else if user.IsOwner() {
		c.Data["json"] = mapping.Mapper.OwnerDbToNet(user.OwnerInfo)
	}
	c.ServeJSONInternal()
}

func (c *ProfileController) PatchMe() {
	admin := &bodies.UpdateProfile{}
	if err := bodies.Require(admin, c.Ctx.Input.RequestBody); err != nil {
		c.BadRequest(err.Error())
	}
	userMe := c.GetUser()
	if admin.Email != nil {
		isOccupied, err := crud.Db.IsEmailOccupied(*admin.Email, &userMe.Id)
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
		isOccupied, err := crud.Db.IsPhoneOccupied(*admin.Phone, &userMe.Id)
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
		adminInfo, err = crud.Db.UpdateAdmin(&crud.UpdateAdminInfo{
			Id:         userMe.Id,
			Surname:    admin.Surname,
			Name:       admin.Name,
			Patronymic: admin.Patronymic,
			Email:      admin.Email,
			Phone:      admin.Phone,
		})
	} else {
		ownerInfo, err = crud.Db.UpdateOwner(&crud.UpdateOwnerInfo{
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
		c.Data["json"] = mapping.Mapper.AdminDbToNet(adminInfo)
	} else if userMe.IsOwner() {
		c.Data["json"] = mapping.Mapper.OwnerDbToNet(ownerInfo)
	} else {
		c.InternalServerError(errors.New("WTF"))
	}
	c.ServeJSONInternal()
}
