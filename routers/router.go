package routers

import (
	"barckend/conf"
	"barckend/controllers"
	"barckend/crud"
	"barckend/security"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	authorizationController := &controllers.AuthorizationController{
		Crud:      crud.Db,
		Encryptor: security.HashMaker,
		TokenManager: security.JWTTokenManager{
			SecretKey: conf.AppConfig.Secret,
		},
	}
	ns := beego.NewNamespace("/api/v1",
		beego.NSRouter("/register", authorizationController, "post:Register"),
		beego.NSRouter("/login", authorizationController, "post:Authorize"),
	)
	beego.AddNamespace(ns)
}
