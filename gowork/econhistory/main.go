// test1 project main.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/axgle/mahonia"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	//"net/url"
	"regexp"
	//"strconv"
	"strings"
	"time"
)

type Body_info struct {
	Results []Econcale `json:"data"`
}
type Econcale struct {
	Date       string `json:"date"`
	Time       string `json:"time"`
	Period     string `json:"period"`
	Importance string `json:"importance"`
	Previous   string `json:"previous"`
	Median     string `json:"median"`
	Ifr_actual string `json:"ifr_actual"`
}
type EconUrl struct {
	Title string
	Url   string
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
}

func timerfunc() {
	fmt.Println("===timer run=====")
}
func starttimer() {
	timer1 := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-timer1.C:
			timerfunc()
		}
	}
}

const dbopen = "root:123456@/wezone?charset=utf8"
const tbname_econdaily = "wz_econdailyfx"
const tbname_econmeeting = "wz_econmeeting"
const tbname_econholiday = "wz_econholiday"

func inittable() {
	//-----数据库操作----
	db, err := sql.Open("mysql", dbopen)
	defer db.Close()
	checkErr(err)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " + tbname_econdaily + " (" +
		"`gid` int(11) DEFAULT NULL," +
		"`date` date DEFAULT NULL," +
		"`time` time DEFAULT NULL," +
		"`currency` varchar(255) DEFAULT NULL," +
		"`country` varchar(255) DEFAULT NULL," +
		"`event` varchar(255) DEFAULT NULL," +
		"`period` varchar(255) DEFAULT NULL," +
		"`importance` varchar(255) DEFAULT NULL," +
		"`previous` varchar(255) DEFAULT NULL," +
		"`median` varchar(255) DEFAULT NULL," +
		"`ifr_actual` varchar(255) DEFAULT NULL," +
		"`gjpy` varchar(255) DEFAULT NULL," +
		"`eventpy` varchar(255) DEFAULT NULL," +
		"PRIMARY KEY (`gid`)" +
		") ENGINE=InnoDB  DEFAULT CHARSET=utf8;")
	checkErr(err)
}
func rega(str_in string) string {
	str_tmp := strings.TrimRight(str_in, ":")
	str_tmp = "\"" + str_tmp + "\":"
	return str_tmp
}
func readdata(str_country string, str_event string, str_url string, db *sql.DB) {
	fmt.Println("查询链接:", str_url)
	resp, errhttp := http.Get(str_url)
	if errhttp != nil {
		fmt.Println("get web fail")
		return
	}
	http_data, _ := ioutil.ReadAll(resp.Body)
	str_data := string(http_data)
	dec := mahonia.NewDecoder("gbk")
	str_utf8 := dec.ConvertString(str_data)
	//fmt.Println(str_utf8)
	//-----数据库操作----
	var str_body string
	str_regexp := "{data:.*}"
	reg := regexp.MustCompile(str_regexp)
	strfind := reg.FindAllString(str_utf8, -1)
	if len(strfind) > 0 {
		str_body = strfind[0]
		//fmt.Println(str_body)
	}
	reg1, _ := regexp.Compile("([a-z_])+:")
	strfind1 := reg1.ReplaceAllStringFunc(str_body, rega)
	//return
	var body Body_info
	if errjson := json.Unmarshal([]byte(strfind1), &body); errjson == nil {
		//======经济数据=====
		for index, result := range body.Results {
			stmt, err1 := db.Prepare("insert " + tbname_econdaily + " set date=?,time=?,country=?,event=?,period=?,importance=?,previous=?,median=?,ifr_actual=?")
			checkErr(err1)
			str_date := result.Date
			str_time := result.Time
			str_period := result.Period
			str_importance := result.Importance
			str_previous := result.Previous
			str_median := result.Median
			str_ifr_actual := result.Ifr_actual
			_, err := stmt.Exec(str_date, str_time, str_country, str_event, str_period, str_importance, str_previous, str_median, str_ifr_actual)
			if err == nil {
				fmt.Printf("====第[%d]组数据开始===\n", index)
				fmt.Printf("时间:%s %s\n", result.Date, result.Time)
				fmt.Printf("国家:%s\n", str_country)
				fmt.Printf("标题:%s\n", str_event)
				fmt.Printf("时期:%s\n", str_period)
				fmt.Printf("前值:%s\n", result.Previous)
				fmt.Printf("预期:%s\n", result.Median)
				fmt.Printf("现值:%s\n", result.Ifr_actual)
				fmt.Printf("====第[%d]组数据结束===\n", index)
			} else {
				fmt.Println(err)
			}
		}
	} else {
		fmt.Println("解析json失败:", errjson)
	}
}
func GetData() {
	inittable()
	//-----数据库操作----
	db, err := sql.Open("mysql", dbopen)
	defer db.Close()
	checkErr(err)
	var urlbase string = "http://money.finance.sina.com.cn/forex/api/jsonp.php/SINAREMOTECALLCALLBACK.CALLBACK_14491916036174448141648899764/FusionChart_Service.getForexChartInfo?"
	var urltime string = "&datefrom=2009-01-01&dateto=2070-11-13&allinfo=1&npp=2000&page=1"
	map_econurl := map[string][]EconUrl{
		"美国": {
			EconUrl{"ISM制造业指数", urlbase + "country=%C3%C0%B9%FA&event=ISM%d6%c6%d4%ec%d2%b5%d6%b8%ca%fd" + urltime},
			EconUrl{"ISM非制造业指数", urlbase + "country=%C3%C0%B9%FA&event=ISM%b7%c7%d6%c6%d4%ec%d2%b5%d6%b8%ca%fd" + urltime},
			EconUrl{"非农就业人数变化", urlbase + "country=%C3%C0%B9%FA&event=%b7%c7%c5%a9%be%cd%d2%b5%c8%cb%ca%fd%b1%e4%bb%af" + urltime},
			EconUrl{"贸易帐", urlbase + "country=%C3%C0%B9%FA&event=%c3%b3%d2%d7%d5%ca" + urltime},
			EconUrl{"失业率", urlbase + "country=%C3%C0%B9%FA&event=%ca%a7%d2%b5%c2%ca" + urltime},
			EconUrl{"未决房屋销售月率", urlbase + "country=%C3%C0%B9%FA&event=%ce%b4%be%f6%b7%bf%ce%dd%cf%fa%ca%db%d4%c2%c2%ca" + urltime},
			EconUrl{"生产者物价指数月率", urlbase + "country=%C3%C0%B9%FA&event=%c9%fa%b2%fa%d5%df%ce%ef%bc%db%d6%b8%ca%fd%d4%c2%c2%ca" + urltime},
			EconUrl{"生产者物价指数年率", urlbase + "country=%C3%C0%B9%FA&event=%c9%fa%b2%fa%d5%df%ce%ef%bc%db%d6%b8%ca%fd%c4%ea%c2%ca" + urltime},
			EconUrl{"核心生产者物价指数月率", urlbase + "country=%C3%C0%B9%FA&event=%ba%cb%d0%c4%c9%fa%b2%fa%d5%df%ce%ef%bc%db%d6%b8%ca%fd%d4%c2%c2%ca" + urltime},
			EconUrl{"核心生产者物价指数年率", urlbase + "country=%C3%C0%B9%FA&event=%ba%cb%d0%c4%c9%fa%b2%fa%d5%df%ce%ef%bc%db%d6%b8%ca%fd%c4%ea%c2%ca" + urltime},
			EconUrl{"核心零售销售月率", urlbase + "country=%C3%C0%B9%FA&event=%ba%cb%d0%c4%c1%e3%ca%db%cf%fa%ca%db%d4%c2%c2%ca" + urltime},
			EconUrl{"零售销售月率", urlbase + "country=%C3%C0%B9%FA&event=%c1%e3%ca%db%cf%fa%ca%db%d4%c2%c2%ca" + urltime},
			EconUrl{"消费者物价指数月率", urlbase + "country=%C3%C0%B9%FA&event=%cf%fb%b7%d1%d5%df%ce%ef%bc%db%d6%b8%ca%fd%d4%c2%c2%ca" + urltime},
			EconUrl{"消费者物价指数年率", urlbase + "country=%C3%C0%B9%FA&event=%cf%fb%b7%d1%d5%df%ce%ef%bc%db%d6%b8%ca%fd%c4%ea%c2%ca" + urltime},
			EconUrl{"核心消费者物价指数月率", urlbase + "country=%C3%C0%B9%FA&event=%ba%cb%d0%c4%cf%fb%b7%d1%d5%df%ce%ef%bc%db%d6%b8%ca%fd%d4%c2%c2%ca" + urltime},
			EconUrl{"核心消费者物价指数年率", urlbase + "country=%C3%C0%B9%FA&event=%ba%cb%d0%c4%cf%fb%b7%d1%d5%df%ce%ef%bc%db%d6%b8%ca%fd%c4%ea%c2%ca" + urltime},
			EconUrl{"新屋开工", urlbase + "country=%C3%C0%B9%FA&event=%d0%c2%ce%dd%bf%aa%b9%a4" + urltime},
			EconUrl{"密歇根大学消费者信心指数初值", urlbase + "country=%C3%C0%B9%FA&event=%c3%dc%d0%aa%b8%f9%b4%f3%d1%a7%cf%fb%b7%d1%d5%df%d0%c5%d0%c4%d6%b8%ca%fd%b3%f5%d6%b5" + urltime},
			EconUrl{"成屋销售", urlbase + "country=%C3%C0%B9%FA&event=%b3%c9%ce%dd%cf%fa%ca%db" + urltime},
			EconUrl{"耐用品订单月率", urlbase + "country=%C3%C0%B9%FA&event=%c4%cd%d3%c3%c6%b7%b6%a9%b5%a5%d4%c2%c2%ca" + urltime},
			EconUrl{"耐用品订单月率（除运输外）", urlbase + "country=%C3%C0%B9%FA&event=%c4%cd%d3%c3%c6%b7%b6%a9%b5%a5%d4%c2%c2%ca%a3%a8%b3%fd%d4%cb%ca%e4%cd%e2%a3%a9" + urltime},
			EconUrl{"咨商会消费者信心指数", urlbase + "country=%C3%C0%B9%FA&event=%d7%c9%c9%cc%bb%e1%cf%fb%b7%d1%d5%df%d0%c5%d0%c4%d6%b8%ca%fd" + urltime},
			EconUrl{"GDP年率初值", urlbase + "country=%C3%C0%B9%FA&event=GDP%c4%ea%c2%ca%b3%f5%d6%b5" + urltime},
			EconUrl{"央行公布利率决议", urlbase + "country=%C3%C0%B9%FA&event=%d1%eb%d0%d0%b9%ab%b2%bc%c0%fb%c2%ca%be%f6%d2%e9" + urltime},
		},
		"德国": {
			EconUrl{"Ifo商业景气指数", urlbase + "country=%b5%c2%b9%fa&event=Ifo%c9%cc%d2%b5%be%b0%c6%f8%d6%b8%ca%fd" + urltime},
			EconUrl{"消费者物价指数月率终值", urlbase + "country=%b5%c2%b9%fa&event=%cf%fb%b7%d1%d5%df%ce%ef%bc%db%d6%b8%ca%fd%d4%c2%c2%ca%d6%d5%d6%b5" + urltime},
			EconUrl{"消费者物价指数年率终值", urlbase + "country=%b5%c2%b9%fa&event=%cf%fb%b7%d1%d5%df%ce%ef%bc%db%d6%b8%ca%fd%c4%ea%c2%ca%d6%d5%d6%b5" + urltime},
			EconUrl{"贸易帐（季调後）", urlbase + "country=%b5%c2%b9%fa&event=%c3%b3%d2%d7%d5%ca%a3%a8%bc%be%b5%f7%e1%e1%a3%a9" + urltime},
			EconUrl{"GDP", urlbase + "country=%b5%c2%b9%fa&event=GDP" + urltime},
			EconUrl{"实际零售销售月率", urlbase + "country=%b5%c2%b9%fa&event=%ca%b5%bc%ca%c1%e3%ca%db%cf%fa%ca%db%d4%c2%c2%ca" + urltime},
			EconUrl{"实际零售销售年率", urlbase + "country=%b5%c2%b9%fa&event=%ca%b5%bc%ca%c1%e3%ca%db%cf%fa%ca%db%c4%ea%c2%ca" + urltime},
			EconUrl{"ZEW经济景气指数", urlbase + "country=%b5%c2%b9%fa&event=ZEW%be%ad%bc%c3%be%b0%c6%f8%d6%b8%ca%fd" + urltime},
		},
		"瑞士": {
			EconUrl{"SVME采购经理人指数", urlbase + "country=%c8%f0%ca%bf&event=SVME%b2%c9%b9%ba%be%ad%c0%ed%c8%cb%d6%b8%ca%fd" + urltime},
			EconUrl{"贸易帐", urlbase + "country=%c8%f0%ca%bf&event=%c3%b3%d2%d7%d5%ca" + urltime},
			EconUrl{"消费者物价指数月率", urlbase + "country=%c8%f0%ca%bf&event=%cf%fb%b7%d1%d5%df%ce%ef%bc%db%d6%b8%ca%fd%d4%c2%c2%ca" + urltime},
			EconUrl{"消费者物价指数年率", urlbase + "country=%c8%f0%ca%bf&event=%cf%fb%b7%d1%d5%df%ce%ef%bc%db%d6%b8%ca%fd%c4%ea%c2%ca" + urltime},
			EconUrl{"GDP季率", urlbase + "country=%c8%f0%ca%bf&event=GDP%bc%be%c2%ca" + urltime},
			EconUrl{"GDP年率", urlbase + "country=%c8%f0%ca%bf&event=GDP%c4%ea%c2%ca" + urltime},
			EconUrl{"央行公布利率决议", urlbase + "country=%c8%f0%ca%bf&event=%d1%eb%d0%d0%b9%ab%b2%bc%c0%fb%c2%ca%be%f6%d2%e9" + urltime},
		},
		"日本": {
			EconUrl{"贸易帐（财务省）", urlbase + "country=%c8%d5%b1%be&event=%c3%b3%d2%d7%d5%ca%a3%a8%b2%c6%ce%f1%ca%a1%a3%a9" + urltime},
			EconUrl{"央行公布利率决议", urlbase + "country=%c8%d5%b1%be&event=%d1%eb%d0%d0%b9%ab%b2%bc%c0%fb%c2%ca%be%f6%d2%e9" + urltime},
			EconUrl{"全国消费者物价指数年率", urlbase + "country=%c8%d5%b1%be&event=%c8%ab%b9%fa%cf%fb%b7%d1%d5%df%ce%ef%bc%db%d6%b8%ca%fd%c4%ea%c2%ca" + urltime},
			EconUrl{"全国核心消费者物价指数年率", urlbase + "country=%c8%d5%b1%be&event=%c8%ab%b9%fa%ba%cb%d0%c4%cf%fb%b7%d1%d5%df%ce%ef%bc%db%d6%b8%ca%fd%c4%ea%c2%ca" + urltime},
			EconUrl{"失业率", urlbase + "country=%c8%d5%b1%be&event=%ca%a7%d2%b5%c2%ca" + urltime},
			EconUrl{"领先指标终值", urlbase + "country=%c8%d5%b1%be&event=%c1%ec%cf%c8%d6%b8%b1%ea%d6%d5%d6%b5" + urltime},
		},
		"英国": {
			EconUrl{"Halifax房价指数月率", urlbase + "country=%d3%a2%b9%fa&event=Halifax%b7%bf%bc%db%d6%b8%ca%fd%d4%c2%c2%ca" + urltime},
			EconUrl{"Halifax房价指数年率", urlbase + "country=%d3%a2%b9%fa&event=Halifax%b7%bf%bc%db%d6%b8%ca%fd%c4%ea%c2%ca" + urltime},
			EconUrl{"贸易帐", urlbase + "country=%d3%a2%b9%fa&event=%c3%b3%d2%d7%d5%ca" + urltime},
			EconUrl{"央行公布利率决议", urlbase + "country=%d3%a2%b9%fa&event=%d1%eb%d0%d0%b9%ab%b2%bc%c0%fb%c2%ca%be%f6%d2%e9" + urltime},
			EconUrl{"核心消费者物价指数年率", urlbase + "country=%d3%a2%b9%fa&event=%ba%cb%d0%c4%cf%fb%b7%d1%d5%df%ce%ef%bc%db%d6%b8%ca%fd%c4%ea%c2%ca" + urltime},
			EconUrl{"核心消费者物价指数月率", urlbase + "country=%d3%a2%b9%fa&event=%ba%cb%d0%c4%cf%fb%b7%d1%d5%df%ce%ef%bc%db%d6%b8%ca%fd%d4%c2%c2%ca" + urltime},
			EconUrl{"消费者物价指数年率", urlbase + "country=%d3%a2%b9%fa&event=%cf%fb%b7%d1%d5%df%ce%ef%bc%db%d6%b8%ca%fd%c4%ea%c2%ca" + urltime},
			EconUrl{"消费者物价指数月率", urlbase + "country=%d3%a2%b9%fa&event=%cf%fb%b7%d1%d5%df%ce%ef%bc%db%d6%b8%ca%fd%d4%c2%c2%ca" + urltime},
			EconUrl{"零售销售月率", urlbase + "country=%d3%a2%b9%fa&event=%c1%e3%ca%db%cf%fa%ca%db%d4%c2%c2%ca" + urltime},
			EconUrl{"零售销售年率", urlbase + "country=%d3%a2%b9%fa&event=%c1%e3%ca%db%cf%fa%ca%db%c4%ea%c2%ca" + urltime},
			EconUrl{"Rightmove房价指数年率", urlbase + "country=%d3%a2%b9%fa&event=Rightmove%b7%bf%bc%db%d6%b8%ca%fd%c4%ea%c2%ca" + urltime},
			EconUrl{"Rightmove房价指数月率", urlbase + "country=%d3%a2%b9%fa&event=Rightmove%b7%bf%bc%db%d6%b8%ca%fd%d4%c2%c2%ca" + urltime},
			EconUrl{"GDP季率初值", urlbase + "country=%d3%a2%b9%fa&event=GDP%bc%be%c2%ca%b3%f5%d6%b5" + urltime},
			EconUrl{"GDP年率初值", urlbase + "country=%d3%a2%b9%fa&event=GDP%c4%ea%c2%ca%b3%f5%d6%b5" + urltime},
			EconUrl{"失业率", urlbase + "country=%d3%a2%b9%fa&event=%ca%a7%d2%b5%c2%ca" + urltime},
		},
		"欧元区": {
			EconUrl{"核心消费者物价指数月率终值", urlbase + "country=%c5%b7%d4%aa%c7%f8&event=%BA%CB%D0%C4%C9%FA%B2%FA%D5%DF%CE%EF%BC%DB%D6%B8%CA%FD%C4%EA%C2%CA" + urltime},
			EconUrl{"核心消费者物价指数年率终值", urlbase + "country=%c5%b7%d4%aa%c7%f8&event=%BA%CB%D0%C4%C9%FA%B2%FA%D5%DF%CE%EF%BC%DB%D6%B8%CA%FD%C4%EA%C2%CA" + urltime},
			EconUrl{"GDP年率终值", urlbase + "country=%c5%b7%d4%aa%c7%f8&event=%BA%CB%D0%C4%C9%FA%B2%FA%D5%DF%CE%EF%BC%DB%D6%B8%CA%FD%C4%EA%C2%CA" + urltime},
			EconUrl{"消费者信心指数终值", urlbase + "country=%c5%b7%d4%aa%c7%f8&event=%BA%CB%D0%C4%C9%FA%B2%FA%D5%DF%CE%EF%BC%DB%D6%B8%CA%FD%C4%EA%C2%CA" + urltime},
			EconUrl{"零售销售月率(%)", urlbase + "country=%c5%b7%d4%aa%c7%f8&event=%BA%CB%D0%C4%C9%FA%B2%FA%D5%DF%CE%EF%BC%DB%D6%B8%CA%FD%C4%EA%C2%CA" + urltime},
			EconUrl{"零售销售年率(%)", urlbase + "country=%c5%b7%d4%aa%c7%f8&event=%BA%CB%D0%C4%C9%FA%B2%FA%D5%DF%CE%EF%BC%DB%D6%B8%CA%FD%C4%EA%C2%CA" + urltime},
			EconUrl{"央行公布利率决议", urlbase + "country=%c5%b7%d4%aa%c7%f8&event=%BA%CB%D0%C4%C9%FA%B2%FA%D5%DF%CE%EF%BC%DB%D6%B8%CA%FD%C4%EA%C2%CA" + urltime},
			EconUrl{"贸易帐（季调後）(亿欧元)", urlbase + "country=%c5%b7%d4%aa%c7%f8&event=%BA%CB%D0%C4%C9%FA%B2%FA%D5%DF%CE%EF%BC%DB%D6%B8%CA%FD%C4%EA%C2%CA" + urltime},
		},
		"澳大利亚": {
			EconUrl{"零售销售月率", urlbase + "country=%b0%c4%b4%f3%c0%fb%d1%c7&event=%c1%e3%ca%db%cf%fa%ca%db%d4%c2%c2%ca" + urltime},
			EconUrl{"贸易帐(亿澳元)", urlbase + "country=%b0%c4%b4%f3%c0%fb%d1%c7&event=%c3%b3%d2%d7%d5%ca(%d2%da%b0%c4%d4%aa)" + urltime},
			EconUrl{"失业率", urlbase + "country=%b0%c4%b4%f3%c0%fb%d1%c7&event=%ca%a7%d2%b5%c2%ca" + urltime},
			EconUrl{"生产者物价指数季率", urlbase + "country=%b0%c4%b4%f3%c0%fb%d1%c7&event=%c9%fa%b2%fa%d5%df%ce%ef%bc%db%d6%b8%ca%fd%bc%be%c2%ca" + urltime},
			EconUrl{"消费者物价指数季率", urlbase + "country=%b0%c4%b4%f3%c0%fb%d1%c7&event=%cf%fb%b7%d1%d5%df%ce%ef%bc%db%d6%b8%ca%fd%bc%be%c2%ca" + urltime},
			EconUrl{"消费者物价指数年率", urlbase + "country=%b0%c4%b4%f3%c0%fb%d1%c7&event=%cf%fb%b7%d1%d5%df%ce%ef%bc%db%d6%b8%ca%fd%c4%ea%c2%ca" + urltime},
			EconUrl{"央行公布利率决议", urlbase + "country=%b0%c4%b4%f3%c0%fb%d1%c7&event=%d1%eb%d0%d0%b9%ab%b2%bc%c0%fb%c2%ca%be%f6%d2%e9" + urltime},
		},
		"加拿大": {
			EconUrl{"新屋开工", urlbase + "country=%bc%d3%c4%c3%b4%f3&event=%d0%c2%ce%dd%bf%aa%b9%a4" + urltime},
			EconUrl{"失业率", urlbase + "country=%bc%d3%c4%c3%b4%f3&event=%ca%a7%d2%b5%c2%ca" + urltime},
			EconUrl{"贸易帐(亿加元)", urlbase + "country=%bc%d3%c4%c3%b4%f3&event=%c3%b3%d2%d7%d5%ca(%d2%da%bc%d3%d4%aa)" + urltime},
			EconUrl{"零售销售月率(%)", urlbase + "country=%bc%d3%c4%c3%b4%f3&event=%c1%e3%ca%db%cf%fa%ca%db%d4%c2%c2%ca(%25)" + urltime},
			EconUrl{"央行公布利率决议", urlbase + "country=%bc%d3%c4%c3%b4%f3&event=%d1%eb%d0%d0%b9%ab%b2%bc%c0%fb%c2%ca%be%f6%d2%e9" + urltime},
			EconUrl{"核心消费者物价指数年率", urlbase + "country=%bc%d3%c4%c3%b4%f3&event=%ba%cb%d0%c4%cf%fb%b7%d1%d5%df%ce%ef%bc%db%d6%b8%ca%fd%c4%ea%c2%ca" + urltime},
			EconUrl{"核心消费者物价指数月率", urlbase + "country=%bc%d3%c4%c3%b4%f3&event=%ba%cb%d0%c4%cf%fb%b7%d1%d5%df%ce%ef%bc%db%d6%b8%ca%fd%d4%c2%c2%ca" + urltime},
			EconUrl{"消费者物价指数年率", urlbase + "country=%bc%d3%c4%c3%b4%f3&event=%cf%fb%b7%d1%d5%df%ce%ef%bc%db%d6%b8%ca%fd%c4%ea%c2%ca" + urltime},
			EconUrl{"消费者物价指数月率", urlbase + "country=%bc%d3%c4%c3%b4%f3&event=%cf%fb%b7%d1%d5%df%ce%ef%bc%db%d6%b8%ca%fd%d4%c2%c2%ca" + urltime},
			EconUrl{"GDP月率(%)", urlbase + "country=%bc%d3%c4%c3%b4%f3&event=GDP%d4%c2%c2%ca(%25)" + urltime},
		},
	}
	//fmt.Println(map_econurl)
	for country, val := range map_econurl {
		for _, econ := range val {
			fmt.Println(country, econ.Title)
			readdata(country, econ.Title, econ.Url, db)
		}
	}
}

func main() {
	GetData()
}
