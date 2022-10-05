package main

import (
	_ "barckend/routers"
	"barckend/setup"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	for _, setuper := range setup.GetSetupers() {
		err := setuper()
		if err != nil {
			panic(err)
		}
	}
	beego.Run()
}
