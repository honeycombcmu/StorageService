package controllers

import (
	"github.com/astaxie/beego"
)

type DashboardController struct {
	beego.Controller
}

func (c *DashboardController) Get() {
	c.TplName = "dashboard.tpl"
	uname := c.GetString("uname")
	c.Data["User_name"] = uname
	c.Data["todo"] = "Please upload your file"
}

func (c *DashboardController) Post() {
	c.TplName = "register.tpl"
	c.SaveToFile("fileField", "uploaded_file.txt")
}
