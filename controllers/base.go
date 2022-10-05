package controllers

import (
	"barckend/crud"
	"barckend/security"
	beego "github.com/beego/beego/v2/server/web"
)

type AuthorizationZone int

const PROFILE_KEY = "payload_key"

const (
	Admin AuthorizationZone = iota
	Owner
)

type BaseController struct {
	beego.Controller
	AuthorizationZones []AuthorizationZone
	TokenManager       security.TokenManager
	Crud               crud.Crud
	Mapper             ModelMapper
}

func (c *BaseController) InternalServerError(err error) {
	c.CustomAbort(500, "Internal server error\n"+err.Error())
}

func (c *BaseController) BadRequest(message string) {
	c.CustomAbort(400, message)
}

func (c *BaseController) Unauthorized() {
	c.CustomAbort(401, "Bad authorization")
}

func (c *BaseController) Prepare() {
	if len(c.AuthorizationZones) > 0 {
		//c.assertAuthorization()
	}
}

//func (c *BaseController) assertAuthorization() {
//	authHeader := c.Ctx.Input.Header("Authorization")
//	if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
//		c.Unauthorized()
//	}
//	bearerAndToken := strings.Split(authHeader, " ")
//	if len(bearerAndToken) != 2 {
//		c.Unauthorized()
//	}
//	payload, err := c.TokenManager.VerifyToken(bearerAndToken[1], security.Access)
//	if err != nil {
//		c.Unauthorized()
//	}
//	profile, err := c.Crud.GetById(payload.UserId)
//	if err != nil {
//		c.InternalServerError(err)
//	}
//	if profile == nil {
//		c.Unauthorized()
//	}
//	c.Data[PROFILE_KEY] = profile
//}
//
//func (c *BaseController) GetProfile() *crud.Profile {
//	profile, exists := c.Data[PROFILE_KEY]
//	if !exists {
//		c.InternalServerError(fmt.Errorf("no profile"))
//	}
//	p, ok := profile.(*crud.Profile)
//	if !ok {
//		c.InternalServerError(fmt.Errorf("profile cast error"))
//	}
//	return p
//}
//
