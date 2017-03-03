// newsinfo project main.go
package main

import (
	"flag"
	// "bufio"
	// "os"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/henrylee2cn/pholcus/app"
	"github.com/henrylee2cn/pholcus/app/spider"
	"github.com/henrylee2cn/pholcus/config"
	//"github.com/henrylee2cn/pholcus/logs"
	"github.com/henrylee2cn/pholcus/runtime/status"
	_ "github.com/henrylee2cn/pholcus/spider_lib" // 此为公开维护的spider规则库
)

var LogicApp = app.New().Init(status.OFFLINE, 0, "")

// 自定义相关配置，将覆盖默认值
func setConf() {
	//mongodb服务器地址
	config.MGO_OUTPUT.Host = "127.0.0.1:27017"
	// mongodb输出时的内容分类
	// key:蜘蛛规则清单
	// value:数据库名
	config.MGO_OUTPUT.DBClass = map[string]string{
		"百度RSS新闻": "1_1",
	}
	// mongodb输出时非默认数据库时以当前时间为集合名
	// h: 精确到小时 (格式 2015-08-28-09)
	// d: 精确到天 (格式 2015-08-28)
	config.MGO_OUTPUT.TableFmt = "d"

	//mysql服务器地址
	config.MYSQL_OUTPUT.Host = "127.0.0.1:3306"
	//msyql数据库
	config.MYSQL_OUTPUT.DefaultDB = "news"
	//mysql用户
	config.MYSQL_OUTPUT.User = "root"
	//mysql密码
	config.MYSQL_OUTPUT.Password = "123456" //
	config.MYSQL_OUTPUT.Table = "`jinshinews`"
}
func Run() {
	// 开启最大核心数运行
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 蜘蛛列表
	var spiderlist string
	for k, v := range LogicApp.GetAllSpiders() {
		spiderlist += "    {" + strconv.Itoa(k) + "} " + v.GetName() + "  " + v.GetDescription() + "\r\n"
	}
	spiderlist = "   【蜘蛛列表】   (选择多蜘蛛以\",\"间隔)\r\n\r\n" + spiderlist
	//spiderflag := flag.String("spider", "", spiderlist+"\r\n")
	// 输出方式
	var outputlib string
	for _, v := range LogicApp.GetOutputLib() {
		outputlib += "{" + v + "} " + v + "    "
	}
	outputlib = strings.TrimRight(outputlib, "    ") + "\r\n"
	outputlib = "   【输出方式】   " + outputlib
	outputflag := flag.String("output", LogicApp.GetOutputLib()[3], outputlib)
	// 并发协程数
	goroutineflag := flag.Uint("go", 20, "   【并发协程】   {1~99999}\r\n")
	// 分批输出
	dockerflag := flag.Uint("docker", 16, "   【分批输出】   每 {1~5000000} 条数据输出一次\r\n")
	// 暂停时间
	pasetimeflag := flag.String("pase", "100,300", "   【暂停时间】   格式如 {基准时间,随机增益} (单位ms)\r\n")
	// 自定义输入
	keywordflag := flag.String("kw", "1", "   【自定义输入<选填>】   多关键词以\",\"隔开\r\n")
	// 采集页数
	maxpageflag := flag.Int("page", 100, "   【采集页数<选填>】\r\n")
	// 备注说明
	flag.String("z", "", "   【说明<非参数>】   各项参数值请参考{}中内容，同一参数包含多个值时以\",\"隔开\r\n\r\n  example：pholcus-cmd.exe -spider=3,8 -output=csv -go=500 -docker=5000 -pase=1000,3000 -kw=pholcus,golang -page=100\r\n")
	flag.Parse()
	// 转换关键词
	keyword := strings.Replace(*keywordflag, ",", "|", -1)
	// 获得暂停时间设置
	var pase [2]uint64
	ptf := strings.Split(*pasetimeflag, ",")
	pase[0], _ = strconv.ParseUint(ptf[0], 10, 64)
	pase[1], _ = strconv.ParseUint(ptf[1], 10, 64)
	// 创建蜘蛛队列
	sps := []*spider.Spider{}
	sps = append(sps, LogicApp.GetSpiderByName("Jinshinew"))
	//if *spiderflag == "" {
	//	logs.Log.Warning(" *     —— 亲，任务列表不能为空哦~")
	//	return
	//}
	//for _, idx := range strings.Split(*spiderflag, ",") {
	//	i, _ := strconv.Atoi(idx)
	//	sps = append(sps, LogicApp.GetAllSpiders()[i])
	//}
	fmt.Println("输出方式:", *outputflag, "\n并发协程数:", *goroutineflag, "\n分批输出数量:", *dockerflag, "\n暂停时间:", *pasetimeflag, "\n关键词:", *keywordflag, "\n采集页数:", *maxpageflag)

	// 配置运行参数
	LogicApp.SetThreadNum(*goroutineflag).
		SetDockerCap(*dockerflag).
		SetOutType(*outputflag).
		SetMaxPage(*maxpageflag).
		SetPausetime([2]uint{uint(pase[0]), uint(pase[1])}).
		SpiderPrepare(sps, keyword).
		Run()
}
func crawl() {
	fmt.Println("===start crawl===")
	if LogicApp.Status() != status.RUN {
		LogicApp.Run()
	}

}
func startTimer() {
	tm := time.NewTicker(30 * time.Minute)
	for {
		select {
		case <-tm.C:
			crawl()
		}
	}
}
func main() {

	setConf() // 不调用则为默认值
	Run()
	startTimer()
}
