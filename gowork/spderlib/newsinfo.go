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
	//"io/ioutil"
	"regexp"
	//"strconv"
	//"strings"

	// 其他包
	"fmt"
	// "math"
	// "time"
)

func init() {
	UPNews.AddMenu()
}

var UPNews = &Spider{
	Name:        "优品财富",
	Description: "优品财富",
	// Pausetime: [2]uint{uint(3000), uint(1000)},
	// Keyword:   USE,
	ResultName: "wz_upnews",
	UseCookie:  false,
	RuleTree: &RuleTree{
		Root: func(self *Spider) {
			fmt.Println("====root====")
			self.AddQueue(map[string]interface{}{
				"Url":  "http://app.upchinafund.com/news/jryjs/columns.html?cid=041001",
				"Rule": "获取列表",
			})
		},

		Trunk: map[string]*Rule{

			"请求列表": {
				ParseFunc: func(self *Spider, resp *context.Response) {
					self.Parse("获取列表", resp)
				},
			},

			"获取列表": {
				ParseFunc: func(self *Spider, resp *context.Response) {
					fmt.Println("====获取列表123====")
					query := resp.GetDom()
					str := resp.GetText()
					fmt.Println(str)
					query.Find(".pic").Each(func(i int, s *goquery.Selection) {
						fmt.Println(i)
						url, _ := s.Find("img").Attr("src")
						fmt.Println(url)
						self.AddQueue(map[string]interface{}{
							"Url":  url,
							"Rule": "输出结果",
						})
					})
				},
			},

			"输出结果": {
				//注意：有无字段语义和是否输出数据必须保持一致
				OutFeild: []string{
					"标题",
					"内容",
					"下载",
				},

				ParseFunc: func(self *Spider, resp *context.Response) {
					query := resp.GetDom()
					fmt.Println("====输出结果====")
					var bt, content, url string
					bt = query.Find(".txt_title").Text()
					content = query.Find(".txt_con").Text()
					//down(72716,'http://img.upchinafund.com/updoc/201511/20104030876.pdf')
					url2, _ := query.Find(".down a").Attr("onclick")
					fmt.Println("====输出结果:====", url2)
					regstr := "http.*pdf"
					reg := regexp.MustCompile(regstr)
					t := reg.FindAllString(url2, -1)
					if len(t) > 0 {
						url = t[0]
						fmt.Println("====输出结果1:====", url)
						self.AddQueue(map[string]interface{}{
							"Url":      url,
							"Rule":     "联系方式",
							"Temp":     map[string]interface{}{"n": bt + ".pdf"},
							"Priority": 1,
						})
					}
					// 结果存入Response中转
					resp.AddItem(map[string]interface{}{
						self.OutFeild(resp, 0): bt,
						self.OutFeild(resp, 1): content,
						self.OutFeild(resp, 2): url,
					})
				},
			},

			"联系方式": {
				ParseFunc: func(self *Spider, resp *context.Response) {
					resp.AddFile(resp.GetTemp("n").(string))
				},
			},
		},
	},
}
