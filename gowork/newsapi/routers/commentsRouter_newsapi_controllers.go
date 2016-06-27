package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["newsapi/controllers:NewsController"] = append(beego.GlobalControllerRouter["newsapi/controllers:NewsController"],
		beego.ControllerComments{
			"GetNews",
			`/getnews`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["newsapi/controllers:NewsController"] = append(beego.GlobalControllerRouter["newsapi/controllers:NewsController"],
		beego.ControllerComments{
			"GetNewsInfo",
			`/getnewsinfo`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["newsapi/controllers:NewsController"] = append(beego.GlobalControllerRouter["newsapi/controllers:NewsController"],
		beego.ControllerComments{
			"GetNewsInfoById",
			`/getnewsinfobyid`,
			[]string{"get"},
			nil})

}
