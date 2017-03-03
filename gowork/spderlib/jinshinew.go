package spider_lib

// 基础包
import (
	//"database/sql"
	//"github.com/PuerkitoBio/goquery" //DOM解析
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
	"encoding/json"
	"fmt"
	// "math"
	// "time"
)

type MainInfo struct {
	Total     int `json:"total"`
	TotalPage int `json:"totalPage"`
	Pre       int `json:"pre"`
}

type ListInfo struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Thumb string `json:"thumb"`
}

type AuthorInfo struct {
	Id   int    `json:"id"`
	Nick string `json:"nick"`
}
type NewsInfo struct {
	Id      int        `json:"id"`
	Title   string     `json:"title"`
	Thumb   string     `json:"thumb"`
	Desc    string     `json:"desc"`
	Time    string     `json:"time_show"`
	Author  AuthorInfo `json:"Author"`
	Content string     `json:"text"`
}

func init() {
	Jinshinew.AddMenu()
}

var Jinshinew = &Spider{
	Name:        "Jinshinew",
	ResultName:  "jinshi",
	Description: "金十资讯实时]",
	// Pausetime: [2]uint{uint(3000), uint(1000)},http://live.wallstreetcn.com/
	Keyword:   USE,
	UseCookie: false,
	RuleTree: &RuleTree{
		Root: func(self *Spider) {
			self.AddQueue(map[string]interface{}{
				"Url":  "https://news.jin10.com/datas/cate/1/main.json",
				"Rule": "请求列表",
				"Temp": map[string]interface{}{"tp": "1", "name": "国际"},
			})
			self.AddQueue(map[string]interface{}{
				"Url":  "https://news.jin10.com/datas/cate/4/main.json",
				"Rule": "请求列表",
				"Temp": map[string]interface{}{"tp": "4", "name": "原油"},
			})
			self.AddQueue(map[string]interface{}{
				"Url":  "https://news.jin10.com/datas/cate/5/main.json",
				"Rule": "请求列表",
				"Temp": map[string]interface{}{"tp": "5", "name": "贵金属"},
			})
			self.AddQueue(map[string]interface{}{
				"Url":  "https://news.jin10.com/datas/cate/6/main.json",
				"Rule": "请求列表",
				"Temp": map[string]interface{}{"tp": "6", "name": "央行"},
			})
			self.AddQueue(map[string]interface{}{
				"Url":  "https://news.jin10.com/datas/cate/7/main.json",
				"Rule": "请求列表",
				"Temp": map[string]interface{}{"tp": "7", "name": "外汇"},
			})
			self.AddQueue(map[string]interface{}{
				"Url":  "https://news.jin10.com/datas/cate/13/main.json",
				"Rule": "请求列表",
				"Temp": map[string]interface{}{"tp": "13", "name": "独家"},
			})
		},

		Trunk: map[string]*Rule{

			"请求列表": {
				ParseFunc: func(self *Spider, resp *context.Response) {
					Body_info := resp.GetText()
					strtype := resp.GetTemp("tp").(string)
					strname := resp.GetTemp("name").(string)
					var maininfo MainInfo
					if errjson := json.Unmarshal([]byte(Body_info), &maininfo); errjson == nil {
						str_totalpage := strconv.Itoa(maininfo.TotalPage)
						strurl := "/p" + str_totalpage + ".json"
						self.AddQueue(map[string]interface{}{
							"Url":  "https://news.jin10.com/datas/cate/" + strtype + strurl,
							"Rule": "文章列表",
							"Temp": map[string]interface{}{
								"type": strtype,
								"name": strname,
							},
						})
					} else {
						fmt.Println("解析json失败:", errjson)
					}

				},
			},

			"文章列表": {
				ParseFunc: func(self *Spider, resp *context.Response) {
					strtype := resp.GetTemp("type").(string)
					strname := resp.GetTemp("name").(string)

					Body_info := resp.GetText()
					var listitems []ListInfo
					if errjson := json.Unmarshal([]byte(Body_info), &listitems); errjson == nil {
						for _, list := range listitems {
							//fmt.Println("=====", list.Id, list.Title)
							self.AddQueue(map[string]interface{}{
								"Url":  "https://news.jin10.com/datas/details/" + strconv.Itoa(list.Id) + ".json",
								"Rule": "详细文章",
								"Temp": map[string]interface{}{
									"type": strtype,
									"name": strname,
								},
							})
						}
					} else {
						fmt.Println("解析List json失败:", errjson)
					}
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
				},
				ParseFunc: func(self *Spider, resp *context.Response) {
					strtype := resp.GetTemp("type").(string)
					strname := resp.GetTemp("name").(string)

					var strarticid, strtitle, strabstract, strtime, strimgname, strimgurl, strautor, strcontent, strcontenthtml string

					Body_info := resp.GetText()
					var newsinfo NewsInfo
					if errjson := json.Unmarshal([]byte(Body_info), &newsinfo); errjson == nil {
						strarticid = strconv.Itoa(newsinfo.Id)
						strtitle = newsinfo.Title
						strabstract = newsinfo.Desc
						strtime = newsinfo.Time
						strimgname = ""
						strimgurl = newsinfo.Thumb
						strautor = newsinfo.Author.Nick
						strcontent = ""
						strcontenthtml = newsinfo.Content
						fmt.Println(strtype, strname, strarticid, strtitle)
					} else {
						fmt.Println("解析newsinfo json失败:", errjson)
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
					})
				},
			},
		},
	},
}
