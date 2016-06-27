package controllers

import (
	//"encoding/json"
	"newsapi/models"

	"github.com/astaxie/beego"
)

// Operations about news
type NewsController struct {
	beego.Controller
}

// @Title getnews
// @Description get news by gid
// @Param	gid		query 	string	false		"The gid for news"
// @Param	limit		query 	string	false		"The litit for news default 30"
// @Success 200 {string} get news success
// @Failure 403 gid not exist
// @router /getnews [get]
func (u *NewsController) GetNews() {
	gid := u.GetString("gid")
	limit := u.GetString("limit")
	jsdata := models.GetNews(gid, limit)
	u.Data["json"] = jsdata
	u.ServeJSON()
}

// @Title getnewsinfo
// @Description get newsinfo by type or time
// @Param	type		query 	string	false		"type:1 国际,type:4 原油,type:5 贵金属,type:6 央行,type:7 外汇,type:13 独家"
// @Param	datetime		query 	string	false		"根据时间来查找"
// @Param	limit		query 	string	false		"The litit for news default 30"
// @Success 200 {string} get news success
// @Failure 403 fail
// @router /getnewsinfo [get]
func (u *NewsController) GetNewsInfo() {
	strtype := u.GetString("type")
	strdatetime := u.GetString("datetime")
	limit := u.GetString("limit")
	jsdata := models.GetNewsInfo(strtype, strdatetime, limit)
	u.Data["json"] = jsdata
	u.ServeJSON()
}

// @Title getnewsinfobyid
// @Description get newsinfo by type or time
// @Param	id		query 	string	false		"文章ID"
// @Success 200 {string} get news success
// @Failure 403 fail
// @router /getnewsinfobyid [get]
func (u *NewsController) GetNewsInfoById() {
	strid := u.GetString("id")
	jsdata := models.GetNewsInfoById(strid)
	u.Data["json"] = jsdata
	u.ServeJSON()
}
