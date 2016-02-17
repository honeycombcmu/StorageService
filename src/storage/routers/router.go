package routers

import (
	"github.com/astaxie/beego"
	"storage/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.MainController{})
	beego.Router("/register", &controllers.UserController{})
	beego.Router("/dashboard", &controllers.DashboardController{})
}
