package controllers

import (
	"barckend/conf"
	"barckend/crud"
	"barckend/security"
	"fmt"
)

type AuthorizationController struct {
	BaseController
	Encryptor security.Hasher
	Mapper    ModelMapper
}

func (c *AuthorizationController) RegisterAdmin() {
	admin := &RegisterAdmin{}
	if err := Require(admin, c.Ctx.Input.RequestBody); err != nil {
		c.BadRequest(err.Error())
	}
	isEmailOccupied, err := c.Crud.IsEmailOccupied(admin.Email, nil)
	if err != nil {
		c.InternalServerError(err)
	}
	if isEmailOccupied {
		c.BadRequest(
			fmt.Sprintf("Profile with this email {%s} is already registered", admin.Email),
		)
	}
	isPhoneOccupied, err := c.Crud.IsPhoneOccupied(admin.Phone, nil)
	if err != nil {
		c.InternalServerError(err)
	}
	if isPhoneOccupied {
		c.BadRequest(
			fmt.Sprintf("Profile with this phone {%s} is already registered", admin.Phone),
		)
	}
	passwordHash, err := security.HashMaker.SHA256(admin.Password)
	if err != nil {
		c.InternalServerError(err)
	}
	adminInfo, err := c.Crud.InsertAdmin(&crud.AdminInfo{
		Name:       admin.Name,
		Surname:    admin.Surname,
		Patronymic: admin.Patronymic,
		Email:      admin.Email,
		Phone:      admin.Phone,
	}, passwordHash)
	if err != nil {
		c.InternalServerError(err)
	}
	c.Data["json"] = c.Mapper.AdminDbToNet(adminInfo)
	if err = c.ServeJSON(); err != nil {
		c.InternalServerError(err)
	}
}

func (c *AuthorizationController) RegisterOwner() {
	owner := &RegisterOwner{}
	if err := Require(owner, c.Ctx.Input.RequestBody); err != nil {
		c.BadRequest(err.Error())
	}
	isEmailOccupied, err := c.Crud.IsEmailOccupied(owner.Email, nil)
	if err != nil {
		c.InternalServerError(err)
	}
	if isEmailOccupied {
		c.BadRequest(
			fmt.Sprintf("Profile with this email {%s} is already registered", owner.Email),
		)
	}
	isPhoneOccupied, err := c.Crud.IsPhoneOccupied(owner.Phone, nil)
	if err != nil {
		c.InternalServerError(err)
	}
	if isPhoneOccupied {
		c.BadRequest(
			fmt.Sprintf("Profile with this phone {%s} is already registered", owner.Phone),
		)
	}
	passwordHash, err := security.HashMaker.SHA256(owner.Password)
	if err != nil {
		c.InternalServerError(err)
	}
	adminInfo, err := c.Crud.InsertOwner(&crud.OwnerInfo{
		Name:       owner.Name,
		Surname:    owner.Surname,
		Patronymic: owner.Patronymic,
		Email:      owner.Email,
		Phone:      owner.Phone,
	}, passwordHash)
	if err != nil {
		c.InternalServerError(err)
	}
	c.Data["json"] = c.Mapper.OwnerDbToNet(adminInfo)
	if err = c.ServeJSON(); err != nil {
		c.InternalServerError(err)
	}
}

func (c *AuthorizationController) Authorize() {
	form := &LoginForm{}
	if err := Require(form, c.Ctx.Input.RequestBody); err != nil {
		c.BadRequest(err.Error())
	}
	passwordHash, err := c.Encryptor.SHA256(form.Password)
	if err != nil {
		c.InternalServerError(err)
	}
	var profile *crud.User
	if form.Phone != nil {
		profile, err = c.Crud.CheckCredentialsPhone(*form.Phone, passwordHash)
	} else {
		profile, err = c.Crud.CheckCredentialsEmail(*form.Email, passwordHash)
	}
	if err != nil {
		c.InternalServerError(err)
	}
	if profile == nil {
		c.BadRequest("Login or password are invalid")
	}
	accessToken, err := c.TokenManager.CreateToken(profile.Id, conf.AppConfig.AccessTokenLifetime, security.Access)
	if err != nil {
		c.InternalServerError(err)
	}
	refreshToken, err := c.TokenManager.CreateToken(profile.Id, conf.AppConfig.RefreshTokenLifetime, security.Refresh)
	if err != nil {
		c.InternalServerError(err)
	}
	c.Data["json"] = &AuthorizationPayload{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Role:         uint(profile.Role),
	}
	if err = c.ServeJSON(); err != nil {
		c.InternalServerError(err)
	}
}
