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
	//"strconv"
	"strings"

	// 其他包
	//"fmt"
	// "math"
	// "time"
)

func init() {
	Jinshi.AddMenu()
}

var Jinshi = &Spider{
	Name:        "金十资讯",
	ResultName:  "jinshi",
	Description: "金十资讯 [http://news.jin10.com/cate/1]",
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
					////把下一个网页放进爬虫里面
					strtype := resp.GetTemp("tp").(string)
					strname := resp.GetTemp("name").(string)
					query := resp.GetDom()
					//开始查找下一页,把下一个网页放进爬虫里面
					s_len := len(query.Find(".pagination li").Nodes)
					if s_len != 0 {
						query.Find(".pagination li").Each(func(i int, s *goquery.Selection) {
							if i == s_len-1 {
								if url, ok := s.Find("a").Attr("href"); ok {
									self.AddQueue(map[string]interface{}{
										"Url":  url,
										"Rule": "请求列表",
										"Temp": map[string]interface{}{"tp": strtype, "name": strname},
									})
								}
							}
						})
					}
				},
			},

			"文章列表": {
				//注意：有无字段语义和是否输出数据必须保持一致
				//OutFeild: []string{
				//	"type",
				//	"name",
				//	"title",
				//	"content",
				//	"time",
				//	"autor",
				//	"imgname",
				//	"imgurl",
				//	"detailurl",
				//},
				ParseFunc: func(self *Spider, resp *context.Response) {
					strtype := resp.GetTemp("tp").(string)
					strname := resp.GetTemp("name").(string)
					query := resp.GetDom()
					var title, content, time, autor, imgurl, detailurl, imgname string
					query.Find(".news_lastest ul li").Each(func(i int, s *goquery.Selection) {
						imgurl, _ = s.Find("img").Attr("src")
						//title, _ = s.Find("img").Attr("alt")
						detailurl, _ = s.Find(".timg a").Attr("href")
						detailurl = "http://news.jin10.com" + detailurl
						title, _ = s.Find(".timg a").Attr("title")
						content = s.Find(".news-p a").Text()
						split := strings.Split(imgurl, "/")
						sz := len(split)
						if sz > 0 {
							imgname = split[sz-1]
							self.AddQueue(map[string]interface{}{
								"Url":  imgurl,
								"Rule": "输出图片",
								"Temp": map[string]interface{}{"n": imgname},
							})
						}
						s.Find(".author_time span").Each(func(i int, s1 *goquery.Selection) {
							if i == 0 {
								autor = s1.Find("a").Text()
							}
							if i == 1 {
								time = s1.Text()
							}
						})

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
							},
						})
						//// 结果存入Response中转
						//resp.AddItem(map[string]interface{}{
						//	self.OutFeild(resp, 0): strtype,
						//	self.OutFeild(resp, 1): strname,
						//	self.OutFeild(resp, 2): title,
						//	self.OutFeild(resp, 3): content,
						//	self.OutFeild(resp, 4): time,
						//	self.OutFeild(resp, 5): autor,
						//	self.OutFeild(resp, 6): imgname,
						//	self.OutFeild(resp, 7): imgurl,
						//	self.OutFeild(resp, 8): detailurl,
						//})
						//////////////////////////////////
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
					strautor := resp.GetTemp("autor").(string)
					strimgname := resp.GetTemp("imgname").(string)
					strimgurl := resp.GetTemp("imgurl").(string)
					strdetailurl := resp.GetTemp("detailurl").(string)
					split := strings.Split(strdetailurl, "/")
					var strarticid string
					if len(split) > 0 {
						strarticid = split[len(split)-1]
					}

					query := resp.GetDom() //Html
					strcontent := query.Find(".newsbox").Text()
					strings.Replace(strcontent, " ", "", -1)
					strings.Replace(strcontent, "金十新闻略作删改", "", -1)
					strings.Replace(strcontent, "金十新闻综合报道", "", -1)
					strings.Replace(strcontent, "金十新闻", "", -1)
					//strcontent := ""
					//strcontenthtml := ""
					strcontenthtml, err := query.Find(".newsbox").Html()
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
