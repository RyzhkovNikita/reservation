package main

import (
	_ "barckend/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.SetStaticPath("api/v1/image/", "images")
	beego.BConfig.WebConfig.DirectoryIndex = true
	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	beego.Run()
}
