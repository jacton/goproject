package spider_lib

// 基础包
import (
	"github.com/PuerkitoBio/goquery"                        //DOM解析
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
	"strings"

	// 其他包
	// "fmt"
	// "math"
	// "time"
)

func init() {
	TianYa.AddMenu()
}

var TianYa = &Spider{
	Name:        "天涯论坛",
	Description: "5万实盘交易日记 [http://bbs.tianya.cn/post-stocks-1107985-1.shtml]",
	// Pausetime: [2]uint{uint(3000), uint(1000)},http://live.wallstreetcn.com/
	Keyword:   USE,
	UseCookie: false,
	RuleTree: &RuleTree{
		Root: func(self *Spider) {
			if strings.IndexAny(self.GetKeyword(), "http://bbs.tianya.cn") == -1 {
				self.SetKeyword("http://bbs.tianya.cn/post-stocks-1107985-1.shtml")
				return
			}
			self.AddQueue(map[string]interface{}{
				"Url":  self.GetKeyword(),
				"Rule": "请求列表",
			})
		},

		Trunk: map[string]*Rule{

			"请求列表": {
				ParseFunc: func(self *Spider, resp *context.Response) {
					// 用指定规则解析响应流
					//把下一个网页放进爬虫里面
					query := resp.GetDom()
					if url, ok := query.Find(".atl-pages .js-keyboard-next").Attr("href"); ok {
						self.AddQueue(map[string]interface{}{
							"Url":  "http://bbs.tianya.cn" + url,
							"Rule": "请求列表",
						})
					}
					self.Parse("输出结果", resp)
				},
			},

			"输出结果": {
				//注意：有无字段语义和是否输出数据必须保持一致
				OutFeild: []string{
					"作者",
					"时间",
					"内容",
					"楼主",
					"楼层",
					"Page",
				},
				ParseFunc: func(self *Spider, resp *context.Response) {
					query := resp.GetDom()
					var 作者, 内容, 时间, 楼主, strpage, 楼层 string
					var bLouzhu bool
					bLouzhu = false
					strpage = query.Find(".atl-head .atl-pages strong").Text()
					nPage, _ := strconv.Atoi(strpage)

					query.Find(".atl-item").Each(func(i int, s *goquery.Selection) {
						作者 = s.Find(".atl-info a").Text()
						内容 = s.Find(".atl-content .bbs-content").Text()
						if i == 0 && nPage == 1 {
							bLouzhu = true
							楼层 = "0楼"
							作者 = query.Find(".atl-menu .atl-info a").Text()
							query.Find(".atl-menu").Find(".atl-info").Find("span").Each(func(i int, sp *goquery.Selection) {
								if i == 1 {
									时间 = sp.Text()
								}
							})
						} else {
							楼层 = s.Find(".atl-content .atl-reply span").Text()
							s.Find(".atl-info").Find("span").Each(func(i int, sp *goquery.Selection) {
								if i == 0 {
									楼主 = s.Find(".host").Text()
									bLouzhu = strings.Contains(楼主, "楼主")
								}
								if i == 1 {
									时间 = sp.Text()
								}
							})
						}

						// 结果存入Response中转
						resp.AddItem(map[string]interface{}{
							self.OutFeild(resp, 0): 作者,
							self.OutFeild(resp, 1): 时间,
							self.OutFeild(resp, 2): 内容,
							self.OutFeild(resp, 3): bLouzhu,
							self.OutFeild(resp, 4): 楼层,
							self.OutFeild(resp, 5): nPage,
						})
						//////////////////////////////////
					})
					////////////////////////
				},
			},
		},
	},
}
