package controllers

import (
	"barckend/conf"
	"barckend/controllers/base"
	"barckend/controllers/requests/bodies"
	"barckend/controllers/responses"
	"barckend/crud"
	"barckend/mapping"
	"barckend/security"
	"fmt"
)

type AuthorizationController struct {
	base.Controller
}

func (c *AuthorizationController) RegisterAdmin() {
	admin := &bodies.RegisterAdmin{}
	if err := bodies.Require(admin, c.Ctx.Input.RequestBody); err != nil {
		c.BadRequest(err.Error())
	}
	isEmailOccupied, err := crud.Db.IsEmailOccupied(admin.Email, nil)
	if err != nil {
		c.InternalServerError(err)
	}
	if isEmailOccupied {
		c.BadRequest(
			fmt.Sprintf("Profile with this email {%s} is already registered", admin.Email),
		)
	}
	isPhoneOccupied, err := crud.Db.IsPhoneOccupied(admin.Phone, nil)
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
	adminInfo, err := crud.Db.InsertAdmin(&crud.AdminInfo{
		Name:       admin.Name,
		Surname:    admin.Surname,
		Patronymic: admin.Patronymic,
		Email:      admin.Email,
		Phone:      admin.Phone,
	}, passwordHash)
	if err != nil {
		c.InternalServerError(err)
	}
	c.ServeJSONInternal(mapping.Mapper.AdminDbToNet(adminInfo))
}

func (c *AuthorizationController) RegisterOwner() {
	owner := &bodies.RegisterOwner{}
	if err := bodies.Require(owner, c.Ctx.Input.RequestBody); err != nil {
		c.BadRequest(err.Error())
	}
	isEmailOccupied, err := crud.Db.IsEmailOccupied(owner.Email, nil)
	if err != nil {
		c.InternalServerError(err)
	}
	if isEmailOccupied {
		c.BadRequest(
			fmt.Sprintf("Profile with this email {%s} is already registered", owner.Email),
		)
	}
	isPhoneOccupied, err := crud.Db.IsPhoneOccupied(owner.Phone, nil)
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
	adminInfo, err := crud.Db.InsertOwner(&crud.OwnerInfo{
		Name:       owner.Name,
		Surname:    owner.Surname,
		Patronymic: owner.Patronymic,
		Email:      owner.Email,
		Phone:      owner.Phone,
	}, passwordHash)
	if err != nil {
		c.InternalServerError(err)
	}
	c.ServeJSONInternal(mapping.Mapper.OwnerDbToNet(adminInfo))
}

func (c *AuthorizationController) Authorize() {
	form := &bodies.LoginForm{}
	if err := bodies.Require(form, c.Ctx.Input.RequestBody); err != nil {
		c.BadRequest(err.Error())
	}
	passwordHash, err := security.HashMaker.SHA256(form.Password)
	if err != nil {
		c.InternalServerError(err)
	}
	var profile *crud.User
	if form.Phone != nil {
		profile, err = crud.Db.CheckCredentialsPhone(*form.Phone, passwordHash)
	} else {
		profile, err = crud.Db.CheckCredentialsEmail(*form.Email, passwordHash)
	}
	if err != nil {
		c.InternalServerError(err)
	}
	if profile == nil {
		c.BadRequest("Login or password are invalid")
	}
	accessToken, err := security.GetTokenManager().CreateToken(profile.Id, conf.AppConfig.AccessTokenLifetime, security.Access)
	if err != nil {
		c.InternalServerError(err)
	}
	refreshToken, err := security.GetTokenManager().CreateToken(profile.Id, conf.AppConfig.RefreshTokenLifetime, security.Refresh)
	if err != nil {
		c.InternalServerError(err)
	}
	c.ServeJSONInternal(&responses.AuthorizationPayload{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Role:         uint(profile.Role),
	})
}
