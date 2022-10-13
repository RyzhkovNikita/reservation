package routers

import (
	"barckend/conf"
	"barckend/controllers"
	"barckend/crud"
	"barckend/security"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	tokenManager := security.JWTTokenManager{
		SecretKey: conf.AppConfig.Secret,
	}
	crudDb := crud.Db
	authorizationController := &controllers.AuthorizationController{
		Encryptor: security.HashMaker,
	}
	authorizationController.Crud = crudDb
	authorizationController.TokenManager = tokenManager
	authorizationController.Mapper = controllers.Mapper

	profileController := &controllers.ProfileController{}
	profileController.Crud = crudDb
	profileController.TokenManager = tokenManager
	profileController.Mapper = controllers.Mapper
	profileController.AuthorizationZones = []crud.Role{crud.Admin, crud.Owner}

	createBarController := &controllers.BarCreateController{}
	createBarController.Crud = crudDb
	createBarController.TokenManager = tokenManager
	createBarController.Mapper = controllers.Mapper
	createBarController.AuthorizationZones = []crud.Role{crud.Owner}

	ns := beego.NewNamespace("/api/v1",
		beego.NSNamespace(
			"/register",
			beego.NSRouter("/admin", authorizationController, "post:RegisterAdmin"),
			beego.NSRouter("/owner", authorizationController, "post:RegisterOwner"),
		),
		beego.NSRouter("/login", authorizationController, "post:Authorize"),
		beego.NSRouter("/me", profileController, "get:GetMe;patch:PatchMe"),
		beego.NSNamespace(
			"/bar",
			beego.NSRouter("/create", createBarController, "post:CreateBar"),
		),
	)
	beego.AddNamespace(ns)
}
