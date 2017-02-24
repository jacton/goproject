package spider_lib

// 基础包
import (
	//"database/sql"
	"github.com/PuerkitoBio/goquery" //DOM解析
	_ "github.com/go-sql-driver/mysql"
	"github.com/henrylee2cn/pholcus/app/downloader/context" //必需
	// "github.com/henrylee2cn/pholcus/logs"               //信息输出
	. "github.com/henrylee2cn/pholcus/app/spider" //必需
	// . "github.com/henrylee2cn/pholcus/app/spider/common"          //选用

	// net包
	// "net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包
	// "regexp"
	"strconv"
	//"strings"

	// 其他包
	"fmt"
	// "math"
	// "time"
)

func init() {
	Gold.AddMenu()
}

var Gold = &Spider{
	Name:        "Gold",
	ResultName:  "gold",
	Description: "gold [http://fx.caiku.com/pair/etf_chart/SPT_GLD/1.html?type=all#etf_data]",
	// Pausetime: [2]uint{uint(3000), uint(1000)},http://live.wallstreetcn.com/
	Keyword:   USE,
	UseCookie: false,
	RuleTree: &RuleTree{
		Root: func(self *Spider) {
			self.AddQueue(map[string]interface{}{
				"Url":  "http://fx.caiku.com/pair/etf_chart/SPT_GLD/1.html?type=all#etf_data",
				"Rule": "请求列表",
				"Temp": map[string]interface{}{"p": 1, "name": "国际"},
			})
		},

		Trunk: map[string]*Rule{

			"请求列表": {
				ParseFunc: func(self *Spider, resp *context.Response) {
					curr := resp.GetTemp("p").(int)
					if c := resp.GetDom().Find("tbody .current").Text(); c != strconv.Itoa(curr) {
						fmt.Println("当前列表页不存在 %v", c)
						return
					}
					self.AddQueue(map[string]interface{}{
						"Url":  "http://fx.caiku.com/pair/etf_chart/SPT_GLD/" + strconv.Itoa(curr+1) + ".html?type=all#etf_data",
						"Rule": "请求列表",
						"Temp": map[string]interface{}{"p": curr + 1},
					})

					// 用指定规则解析响应流
					self.Parse("详细信息", resp)
				},
			},

			"详细信息": {
				OutFeild: []string{
					"时间",
					"持仓量(盎司)",
					"持仓量(吨)",
					"总价值(美元)",
					"持仓量增减(盎司)",
				},
				ParseFunc: func(self *Spider, resp *context.Response) {
					query := resp.GetDom()
					var date, chicang_as, chicang_dun, alljiazhi, chicangliang string
					query.Find(".tab01 tbody tr").Each(func(i int, s *goquery.Selection) {
						//strcontent := s.Text()
						s.Find("td").Each(func(i int, s1 *goquery.Selection) {
							switch i {
							case 0:
								date = s1.Text()
								break
							case 1:
								chicang_as = s1.Text()
								break
							case 2:
								chicang_dun = s1.Text()
								break
							case 3:
								alljiazhi = s1.Text()
								break
							case 4:
								chicangliang = s1.Text()
								break
							}
						})
						fmt.Println("==========")
						fmt.Println(date)
						//return
						// 结果存入Response中转
						resp.AddItem(map[string]interface{}{
							self.OutFeild(resp, 0): date,
							self.OutFeild(resp, 1): chicang_as,
							self.OutFeild(resp, 2): chicang_dun,
							self.OutFeild(resp, 3): alljiazhi,
							self.OutFeild(resp, 4): chicangliang,
						})
					})
					////////////////////////
				},
			},
			///////////////////////////////////////////////////
		},
	},
}
