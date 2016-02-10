package routers

import (
	"storage/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.MainController{})
	beego.Router("/register", &controllers.UserController{})
}
