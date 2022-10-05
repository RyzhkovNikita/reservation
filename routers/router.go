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

	ns := beego.NewNamespace("/api/v1",
		beego.NSNamespace(
			"/register",
			beego.NSRouter("/admin", authorizationController, "post:RegisterAdmin"),
			beego.NSRouter("/owner", authorizationController, "post:RegisterOwner"),
		),
		beego.NSRouter("/login", authorizationController, "post:Authorize"),
		//beego.NSNamespace(
		//	"/profile",
		//	beego.NSRouter("", profileController, "get:GetBarProfile"),
		//	beego.NSRouter("", profileController, "patch:UpdateProfile"),
		//	beego.NSRouter("/logo", profileController, "post:UploadLogo"),
		//),
	)
	beego.AddNamespace(ns)
}
