package routers

import (
	"barckend/controllers"
	"barckend/controllers/bar"
	"barckend/crud"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	authorization := &controllers.AuthorizationController{}

	profile := &controllers.ProfileController{}
	profile.AuthorizationZones = []crud.Role{crud.Admin, crud.Owner}
	createBar := &bar.CreateController{}
	createBar.AuthorizationZones = []crud.Role{crud.Owner}
	getBarInfo := &bar.InfoController{}
	getBarInfo.AuthorizationZones = []crud.Role{crud.Owner, crud.Admin}
	editBarInfo := &bar.EditController{}
	editBarInfo.AuthorizationZones = []crud.Role{crud.Owner}

	ns := beego.NewNamespace("/api/v1",
		beego.NSNamespace(
			"/register",
			beego.NSRouter("/admin", authorization, "post:RegisterAdmin"),
			beego.NSRouter("/owner", authorization, "post:RegisterOwner"),
		),
		beego.NSRouter("/login", authorization, "post:Authorize"),
		beego.NSRouter("/me", profile, "get:GetMe;patch:PatchMe"),
		beego.NSNamespace(
			"/bar",
			beego.NSRouter("/create", createBar, "post:CreateBar"),
			beego.NSRouter("/:bar_id([0-9]+)", getBarInfo, "get:GetBarInformation"),
			beego.NSRouter("/:bar_id([0-9]+)", editBarInfo, "patch:EditBar"),
		),
	)
	beego.AddNamespace(ns)
}
