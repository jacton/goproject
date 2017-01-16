package spider_lib

// 基础包
import (
	"database/sql"
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
	//"strconv"
	"strings"

	// 其他包
	"fmt"
	// "math"
	// "time"
)

func init() {
	Jinshi1.AddMenu()
}

var db *sql.DB
var Jinshi1 = &Spider{
	Name:        "金十资讯实时",
	ResultName:  "jinshi",
	Description: "金十资讯实时 [http://news.jin10.com/cate/1]",
	// Pausetime: [2]uint{uint(3000), uint(1000)},http://live.wallstreetcn.com/
	Keyword:   USE,
	UseCookie: false,
	RuleTree: &RuleTree{
		Root: func(self *Spider) {
			self.AddQueue(map[string]interface{}{
				"Url":  "http://news.jin10.com/cate/1",
				"Rule": "请求列表",
				"Temp": map[string]interface{}{"tp": "1", "name": "国际"},
			})
			self.AddQueue(map[string]interface{}{
				"Url":  "http://news.jin10.com/cate/4",
				"Rule": "请求列表",
				"Temp": map[string]interface{}{"tp": "4", "name": "原油"},
			})
			self.AddQueue(map[string]interface{}{
				"Url":  "http://news.jin10.com/cate/5",
				"Rule": "请求列表",
				"Temp": map[string]interface{}{"tp": "5", "name": "贵金属"},
			})
			self.AddQueue(map[string]interface{}{
				"Url":  "http://news.jin10.com/cate/6",
				"Rule": "请求列表",
				"Temp": map[string]interface{}{"tp": "6", "name": "央行"},
			})
			self.AddQueue(map[string]interface{}{
				"Url":  "http://news.jin10.com/cate/7",
				"Rule": "请求列表",
				"Temp": map[string]interface{}{"tp": "7", "name": "外汇"},
			})
			self.AddQueue(map[string]interface{}{
				"Url":  "http://news.jin10.com/cate/13",
				"Rule": "请求列表",
				"Temp": map[string]interface{}{"tp": "13", "name": "独家"},
			})
		},

		Trunk: map[string]*Rule{

			"请求列表": {
				ParseFunc: func(self *Spider, resp *context.Response) {
					// 用指定规则解析每篇文章
					self.Parse("文章列表", resp)
				},
			},

			"文章列表": {
				ParseFunc: func(self *Spider, resp *context.Response) {
					strtype := resp.GetTemp("tp").(string)
					strname := resp.GetTemp("name").(string)
					query := resp.GetDom()
					if db == nil {
						db, _ = sql.Open("mysql", "root:123456@/pholcus?charset=utf8")
					}
					var title, content, time, autor, imgurl, detailurl, imgname string
					query.Find(".jin-newsList__item").Each(func(i int, s *goquery.Selection) {
						imgurl, _ = s.Find(".jin-newsList__img img").Attr("data-original")
						//imgurl, _ = s.Find(".jin-newsList__img").Html()
						//title, _ = s.Find("img").Attr("alt")
						detailurl, _ = s.Find("a").Attr("href")
						detailurl = "http://news.jin10.com" + detailurl
						title = s.Find(".jin-newsList__title").Text()
						content = title
						//split := strings.Split(imgurl, "/")
						//sz := len(split)
						//if sz > 0 {
						//	imgname = split[sz-1]
						//	self.AddQueue(map[string]interface{}{
						//		"Url":  imgurl,
						//		"Rule": "输出图片",
						//		"Temp": map[string]interface{}{"n": imgname},
						//	})
						//}
						//tm := s.Find(".jin-newsList__time").Text()
						//dt := s.Find(".jin-newsList__icon-float").Text()
						//time = dt + " " + tm + ":00"
						fmt.Println("==========")
						fmt.Println(title)
						//fmt.Println(time)

						split := strings.Split(detailurl, "/")
						var strarticid string
						if len(split) > 0 {
							strarticid = split[len(split)-1]
						}
						//////////////////////业务逻辑处理
						var articid = ""
						db.QueryRow("SELECT `articid` FROM `金十资讯-详细文章-1` WHERE `type`=? and `articid` = ?", strtype, strarticid).Scan(&articid)
						if articid == "" {
							//////////////////////////////
							self.AddQueue(map[string]interface{}{
								"Url":  detailurl,
								"Rule": "详细文章",
								"Temp": map[string]interface{}{
									"type":      strtype,
									"name":      strname,
									"title":     title,
									"abstract":  content,
									"time":      time,
									"autor":     autor,
									"imgname":   imgname,
									"imgurl":    imgurl,
									"detailurl": detailurl,
									"articid":   strarticid,
								},
							})
							////////////////////////////////////
						}
					})
					////////////////////////
				},
			},
			"输出图片": {
				ParseFunc: func(self *Spider, resp *context.Response) {
					resp.AddFile(resp.GetTemp("n").(string))
				},
			},
			"详细文章": {
				OutFeild: []string{
					"type",
					"name",
					"articid",
					"title",
					"abstract",
					"content",
					"contenthtml",
					"time",
					"autor",
					"imgname",
					"imgurl",
					"detailurl",
				},
				ParseFunc: func(self *Spider, resp *context.Response) {
					strtype := resp.GetTemp("type").(string)
					strname := resp.GetTemp("name").(string)
					strtitle := resp.GetTemp("title").(string)
					strabstract := resp.GetTemp("abstract").(string)
					strtime := resp.GetTemp("time").(string)
					//strautor := resp.GetTemp("autor").(string)
					strimgname := resp.GetTemp("imgname").(string)
					strimgurl := resp.GetTemp("imgurl").(string)
					strdetailurl := resp.GetTemp("detailurl").(string)
					split := strings.Split(strdetailurl, "/")
					var strarticid string
					if len(split) > 0 {
						strarticid = split[len(split)-1]
					}

					query := resp.GetDom() //Html
					var strautor string

					query.Find(".jin-meta span").Each(func(i int, s1 *goquery.Selection) {
						if i == 2 {
							strtime = s1.Text()
						}
						if i == 3 {
							strautor = s1.Text()
						}
					})
					strcontent := query.Find(".jin-news-article_content").Text()
					strings.Replace(strcontent, " ", "", -1)
					//strcontent := ""
					//strcontenthtml := ""
					strcontenthtml, err := query.Find(".jin-news-article_content").Html()
					strings.Replace(strcontenthtml, " ", "", -1)
					if err != nil {
						return
					}
					// 结果存入Response中转
					resp.AddItem(map[string]interface{}{
						self.OutFeild(resp, 0):  strtype,
						self.OutFeild(resp, 1):  strname,
						self.OutFeild(resp, 2):  strarticid,
						self.OutFeild(resp, 3):  strtitle,
						self.OutFeild(resp, 4):  strabstract,
						self.OutFeild(resp, 5):  strcontent,
						self.OutFeild(resp, 6):  strcontenthtml,
						self.OutFeild(resp, 7):  strtime,
						self.OutFeild(resp, 8):  strautor,
						self.OutFeild(resp, 9):  strimgname,
						self.OutFeild(resp, 10): strimgurl,
						self.OutFeild(resp, 11): strdetailurl,
					})
				},
			},
		},
	},
}
