package controllers

import (
	"barckend/crud"
	"barckend/security"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"golang.org/x/exp/slices"
	"strings"
)

type AuthorizationZone int

const PROFILE_KEY = "payload_key"

type BaseController struct {
	beego.Controller
	AuthorizationZones []crud.Role
	TokenManager       security.TokenManager
	Crud               crud.AdminCrud
}

func (c *BaseController) InternalServerError(err error) {
	c.CustomAbort(500, "Internal server error\n"+err.Error())
}

func (c *BaseController) BadRequest(message string) {
	c.CustomAbort(400, message)
}

func (c *BaseController) Unauthorized() {
	c.CustomAbort(401, "")
}

func (c *BaseController) Forbidden() {
	c.CustomAbort(403, "")
}

func (c *BaseController) ServeJSONInternal() {
	if err := c.ServeJSON(); err != nil {
		c.InternalServerError(err)
	}
}

func (c *BaseController) Prepare() {
	if len(c.AuthorizationZones) > 0 {
		c.assertAuthorization()
	}
}

func (c *BaseController) assertAuthorization() {
	authHeader := c.Ctx.Input.Header("Authorization")
	if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		c.Unauthorized()
	}
	bearerAndToken := strings.Split(authHeader, " ")
	if len(bearerAndToken) != 2 {
		c.Unauthorized()
	}
	payload, err := c.TokenManager.VerifyToken(bearerAndToken[1], security.Access)
	if err != nil {
		c.Unauthorized()
	}
	profile, err := c.Crud.GetById(payload.UserId)
	if err != nil {
		c.InternalServerError(err)
	}
	if profile == nil {
		c.Unauthorized()
	}
	if !slices.Contains(c.AuthorizationZones, profile.Role) {
		c.Forbidden()
	}
	c.Data[PROFILE_KEY] = profile
}

func (c *BaseController) GetUser() *crud.User {
	user, exists := c.Data[PROFILE_KEY]
	if !exists {
		c.InternalServerError(fmt.Errorf("no user"))
	}
	p, ok := user.(*crud.User)
	if !ok {
		c.InternalServerError(fmt.Errorf("user cast error"))
	}
	return p
}
