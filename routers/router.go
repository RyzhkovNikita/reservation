package routers

import (
	"barckend/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/api/v1/register", &controllers.RegistrationController{})
}
