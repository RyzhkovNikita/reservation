package main

import (
	_ "barckend/routers"
	"barckend/setup"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

func main() {
	for _, setuper := range setup.GetSetupers() {
		err := setuper()
		if err != nil {
			panic(err)
		}
	}
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowCredentials: true,
	}))
	beego.Run()
}
