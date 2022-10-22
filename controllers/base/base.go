package base

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

type Controller struct {
	beego.Controller
	AuthorizationZones []crud.Role
}

func (c *Controller) InternalServerError(err error) {
	c.CustomAbort(500, "Internal server error\n"+err.Error())
}

func (c *Controller) BadRequest(message string) {
	c.CustomAbort(400, message)
}

func (c *Controller) Unauthorized() {
	c.CustomAbort(401, "")
}

func (c *Controller) Forbidden() {
	c.CustomAbort(403, "")
}

func (c *Controller) NotFound(message string) {
	c.CustomAbort(404, message)
}

func (c *Controller) ServeJSONInternal() {
	if err := c.ServeJSON(); err != nil {
		c.InternalServerError(err)
	}
}

func (c *Controller) Prepare() {
	if len(c.AuthorizationZones) > 0 {
		c.assertAuthorization()
	}
}

func (c *Controller) assertAuthorization() {
	authHeader := c.Ctx.Input.Header("Authorization")
	if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		c.Unauthorized()
	}
	bearerAndToken := strings.Split(authHeader, " ")
	if len(bearerAndToken) != 2 {
		c.Unauthorized()
	}
	payload, err := security.GetTokenManager().VerifyToken(bearerAndToken[1], security.Access)
	if err != nil {
		c.Unauthorized()
	}
	profile, err := crud.Db.GetById(payload.UserId)
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

func (c *Controller) GetUser() *crud.User {
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
