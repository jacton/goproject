package routers

import (
	"github.com/astaxie/beego"
	"newsweb/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/newsinfo", &controllers.NewsController{})
}
