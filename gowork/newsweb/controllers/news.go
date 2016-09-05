package controllers

import (
	"github.com/astaxie/beego"
	"newsweb/models"
)

type NewsController struct {
	beego.Controller
}

func (c *NewsController) Get() {
	strid := c.GetString("id")
	jsdata := models.GetNewsInfoById(strid)
	if jsdata["retcode"] == 0 {
		var tmpinterface interface{}
		tmpinterface = jsdata["data"]
		var tmpNewsInfo models.NewsInfo = tmpinterface.(models.NewsInfo)

		c.Data["NewsTitle"] = tmpNewsInfo.Title
		c.Data["NewsTypeName"] = tmpNewsInfo.Name
		c.Data["NewsAutor"] = tmpNewsInfo.Autor
		c.Data["NewsSrc"] = tmpNewsInfo.Autor
		c.Data["NewsDateTime"] = tmpNewsInfo.Datetime
		c.Data["NewsContent"] = tmpNewsInfo.Contenthtml
		c.Data["NewsImg"] = tmpNewsInfo.Imgurl
		c.TplName = "news.tpl"
	} else if jsdata["retcode"] == 1 {
		c.TplName = "404.tpl"
	} else {
		c.Abort("404")
	}

}
