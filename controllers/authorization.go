package controllers

import (
	"barckend/conf"
	"barckend/crud"
	"barckend/security"
	"barckend/utils"
	"fmt"
)

type AuthorizationController struct {
	utils.BaseController
	Crud         crud.Crud
	Encryptor    security.Hasher
	TokenManager security.TokenManager
}

func (c *AuthorizationController) Register() {
	bar := &RegisterBar{}
	if err := bar.ParseFromBody(c.Ctx.Input.RequestBody); err != nil {
		c.BadRequest(err.Error())
	}
	existingProfile, err := c.Crud.GetByEmail(bar.Email)
	if err != nil {
		c.InternalServerError(err)
	}
	if existingProfile != nil {
		c.BadRequest(
			fmt.Sprintf("Profile with this email {%s} is already existing", bar.Email),
		)
	}
	passwordHash, err := security.HashMaker.CalculateHash(bar.Password)
	if err != nil {
		c.InternalServerError(err)
	}
	profile, err := c.Crud.Insert(&crud.Profile{
		Name:        bar.Name,
		Description: bar.Description,
		Address:     bar.Address,
		LogoUrl:     "",
		Credentials: &crud.Credentials{
			Email:        bar.Email,
			PasswordHash: passwordHash,
		},
	})
	if err != nil {
		c.InternalServerError(err)
	}
	c.Data["json"] = &BarInfoResponse{
		Id:          profile.Id,
		Email:       profile.Credentials.Email,
		Name:        profile.Name,
		Description: profile.Description,
		Address:     profile.Address,
		LogoUrl:     nil,
	}
	if err = c.ServeJSON(); err != nil {
		c.InternalServerError(err)
	}
}

func (c *AuthorizationController) Authorize() {
	form := &LoginForm{}
	if err := form.ParseFromBody(c.Ctx.Input.RequestBody); err != nil {
		c.BadRequest(err.Error())
	}
	passwordHash, err := c.Encryptor.CalculateHash(form.Password)
	if err != nil {
		c.InternalServerError(err)
	}
	profile, err := c.Crud.CheckCredentials(form.Email, passwordHash)
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
	}
	if err = c.ServeJSON(); err != nil {
		c.InternalServerError(err)
	}
}
