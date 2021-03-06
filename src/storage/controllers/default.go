package controllers

import (
	"github.com/astaxie/beego"
	"log"
	"storage/models"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	beego.SetStaticPath("/static", "static")
	beego.SetStaticPath("/images", "images")
	beego.SetStaticPath("/css", "css")
	beego.SetStaticPath("/js", "js")

	c.TplName = "login.tpl"
}

func (c *MainController) Post() {
	// Check validation of user/pass
	//c.Ctx.WriteString("Your user name is: xxx")
	c.TplName = "login.tpl"
	uname := c.GetString("email")
	password := c.GetString("password")
	unameKey := models.FormatUserLoginKey(uname)
	storedPwd, err := models.DB.Get(unameKey)
	if err != nil {
		log.Println("user:" + uname + " does not exist.")
		return // user does not exist, just return password and username don't match
	}
	if models.Hash(password) != storedPwd {
		log.Println("password is wrong")
		return // password is wrong, just return password and username don't match
	}
	c.Redirect("/dashboard?uname="+uname, 302)
}
