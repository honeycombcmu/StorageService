package controllers

import (
	"github.com/astaxie/beego"
    "storage/models"
	"log"
)

type UserController struct {
	beego.Controller
}

func (c *UserController) Get() {
	c.TplNames = "register.tpl"
}

func (c *UserController) Post() {
	c.TplNames = "register.tpl"
	uname := c.GetString("email")
	password := c.GetString("password")
	rePassword := c.GetString("re-password")
	if password != rePassword {
		return
	}
	unameKey := models.FormatUserLoginKey(uname)
	_, err := models.DB.Get(unameKey)
	if err == nil {
		log.Println("user: " + uname + " exists.")
		return // user exists
	}

    err = models.DB.Put(unameKey, models.Hash(password))
    if err != nil {
        log.Println("failed to put into database")
        return // failure
    }

	c.Redirect("/login?user_name="+uname, 302)
}
