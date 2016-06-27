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
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Body_info struct {
	Results []Econcale    `json:"data"`
	Meeting []Econmeeting `json:"meeting"`
	Holiday []Econholiday `json:"holiday"`
}
type Econcale struct {
	Id         string `json:"id"`
	Date       string `json:"date"`
	Time       string `json:"time"`
	Currency   string `json:"currency"`
	Country    string `json:"country"`
	Event      string `json:"event"`
	Period     string `json:"period"`
	Importance string `json:"importance"`
	Previous   string `json:"previous"`
	Median     string `json:"median"`
	Ifr_actual string `json:"ifr_actual"`
	Gjpy       string `json:"gjpy"`
	Eventpy    string `json:"eventpy"`
}
type Econmeeting struct {
	Id          string `json:"id"`
	Currency    string `json:"currency"`
	Title       string `json:"title"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	Time_hidden string `json:"time_hidden"`
	Detail      string `json:"detail"`
	Insert_date string `json:"insert_date"`
}
type Econholiday struct {
	Id          string `json:"id"`
	Currency    string `json:"currency"`
	Title       string `json:"title"`
	Date        string `json:"date"`
	Insert_date string `json:"insert_date"`
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		//panic(err)
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
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " + tbname_econmeeting + " (" +
		"`gid` int(11) NOT NULL," +
		"`currency` varchar(255) DEFAULT NULL," +
		"`title` varchar(1000) DEFAULT NULL," +
		"`date` date DEFAULT NULL," +
		"`time` time DEFAULT NULL," +
		"`time_hidden` time DEFAULT NULL," +
		"`detail` text," +
		"`insert_date` datetime DEFAULT NULL," +
		"PRIMARY KEY (`gid`)" +
		") ENGINE=InnoDB  DEFAULT CHARSET=utf8;")
	checkErr(err)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " + tbname_econholiday + " (" +
		"`gid` int(11) NOT NULL," +
		"`currency` varchar(255) DEFAULT NULL," +
		"`title` varchar(255) DEFAULT NULL," +
		"`date` date DEFAULT NULL," +
		"`insert_date` datetime DEFAULT NULL," +
		"PRIMARY KEY (`gid`)" +
		") ENGINE=InnoDB  DEFAULT CHARSET=utf8;")
	checkErr(err)
	//-----------------
}
func rega(str_in string) string {
	str_tmp := strings.TrimRight(str_in, ":")
	str_tmp = "\"" + str_tmp + "\":"
	return str_tmp
}
func updata(str_date string, db *sql.DB) { //更新数据
	//str_time="2015-11-27"
	fmt.Println("===start up data====")
	fmt.Println("查询的日期:", str_date)
	str_baseurl := "http://vip.stock.finance.sina.com.cn/forex/api/jsonp.php/SINAFINANCE144893533921925307/DailyFX_AllService.getOnedayEventsNew?date="
	str_baseurl = str_baseurl + str_date
	fmt.Println("查询链接:", str_baseurl)
	resp, errhttp := http.Get(str_baseurl)
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
	//db, err := sql.Open("mysql", dbopen)
	//defer db.Close()
	//checkErr(err)
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
			var n_gid int
			n_gid, _ = strconv.Atoi(result.Id)
			str_ifr_actual := result.Ifr_actual
			if len(str_ifr_actual) > 0 {
				stmt, err1 := db.Prepare("update " + tbname_econdaily + " set ifr_actual=? where gid=?")
				checkErr(err1)
				_, err := stmt.Exec(str_ifr_actual, n_gid)
				checkErr(err)
				fmt.Printf("====第[%d]组数据开始===\n", index)
				fmt.Printf("id:%s\n", result.Id)
				fmt.Printf("时间:%s %s\n", result.Date, result.Time)
				fmt.Printf("标题:%s\n", result.Event)
				fmt.Printf("国家:%s\n", result.Country)
				fmt.Printf("前值:%s\n", result.Previous)
				fmt.Printf("预期:%s\n", result.Median)
				fmt.Printf("现值:%s\n", result.Ifr_actual)
				fmt.Printf("====第[%d]组数据结束===\n", index)
			}
		}
	} else {
		fmt.Println("解析json失败:", errjson)
	}
}
func updatarow(str_date string, str_time string, db *sql.DB) { //更新数据
	//str_time="2015-11-27"
	fmt.Println("===start up data====")
	fmt.Println("查询的日期:", str_date)
	str_baseurl := "http://vip.stock.finance.sina.com.cn/forex/api/jsonp.php/SINAFINANCE144893533921925307/DailyFX_AllService.getOnedayEventsNew?date="
	str_baseurl = str_baseurl + str_date
	fmt.Println("查询链接:", str_baseurl)
	resp, errhttp := http.Get(str_baseurl)
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
	}
	reg1, _ := regexp.Compile("([a-z_])+:")
	strfind1 := reg1.ReplaceAllStringFunc(str_body, rega)
	//return
	var body Body_info
	if errjson := json.Unmarshal([]byte(strfind1), &body); errjson == nil {
		//======经济数据=====
		tm1, _ := time.Parse("15:04:05", str_time) //传入的当天没有数据的最小值
		for _, result := range body.Results {
			tm2, _ := time.Parse("15:04:05", result.Time)
			//fmt.Println("===time1===", tm1)
			//fmt.Println("===time2===", tm2)
			if tm2.After(tm1) == true || tm2.Equal(tm1) {
				//fmt.Printf("时间:%s %s\n", result.Date, result.Time)
				//fmt.Printf("标题:%s\n", result.Event)
				//fmt.Printf("现值:%s\n", result.Ifr_actual)
				var n_gid int
				n_gid, _ = strconv.Atoi(result.Id)
				str_ifr_actual := result.Ifr_actual
				if len(str_ifr_actual) > 0 {
					stmt, err1 := db.Prepare("update " + tbname_econdaily + " set ifr_actual=? where gid=?")
					checkErr(err1)
					_, err := stmt.Exec(str_ifr_actual, n_gid)
					if err == nil {
						fmt.Println("更新数据成功!")
						fmt.Printf("id:%s\n", result.Id)
						fmt.Printf("时间:%s %s\n", result.Date, result.Time)
						fmt.Printf("标题:%s\n", result.Event)
						//fmt.Printf("国家:%s\n", result.Country)
						//fmt.Printf("前值:%s\n", result.Previous)
						//fmt.Printf("预期:%s\n", result.Median)
						//fmt.Printf("现值:%s\n", result.Ifr_actual)
					}
				}
			}
		}
	} else {
		fmt.Println("解析json失败:", errjson)
	}
}

func readdata(str_time string, db *sql.DB) {
	//str_time="2015-11-27"
	fmt.Println("===start get data====")
	fmt.Println("查询的日期:", str_time)
	str_baseurl := "http://vip.stock.finance.sina.com.cn/forex/api/jsonp.php/SINAFINANCE144893533921925307/DailyFX_AllService.getOnedayEventsNew?date="
	str_baseurl = str_baseurl + str_time
	fmt.Println("查询链接:", str_baseurl)
	resp, errhttp := http.Get(str_baseurl)
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
	//db, err := sql.Open("mysql",dbopen)
	//defer db.Close()
	//checkErr(err)
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
			stmt, err1 := db.Prepare("insert " + tbname_econdaily + " set gid=?,date=?,time=?,currency=?,country=?,event=?,period=?,importance=?,previous=?,median=?,ifr_actual=?,gjpy=?,eventpy=?")
			checkErr(err1)
			var n_gid int
			n_gid, _ = strconv.Atoi(result.Id)
			str_date := result.Date
			str_time := result.Time
			str_currency := result.Currency
			str_country := result.Country
			str_event := result.Event
			str_period := result.Period
			str_importance := result.Importance
			str_previous := result.Previous
			str_median := result.Median
			str_ifr_actual := result.Ifr_actual
			str_gjpy := result.Gjpy
			str_eventpy := result.Eventpy
			_, err := stmt.Exec(n_gid, str_date, str_time, str_currency, str_country, str_event, str_period, str_importance, str_previous, str_median, str_ifr_actual, str_gjpy, str_eventpy)
			if err == nil {
				fmt.Printf("====第[%d]组数据开始===\n", index)
				fmt.Printf("id:%s\n", result.Id)
				fmt.Printf("时间:%s %s\n", result.Date, result.Time)
				fmt.Printf("标题:%s\n", result.Event)
				fmt.Printf("国家:%s\n", result.Country)
				fmt.Printf("前值:%s\n", result.Previous)
				fmt.Printf("预期:%s\n", result.Median)
				fmt.Printf("现值:%s\n", result.Ifr_actual)
				fmt.Printf("====第[%d]组数据结束===\n", index)
			} else {
				fmt.Println(err)
			}
		}
		//=====财经大事=======
		for index, result := range body.Meeting {
			stmt, err1 := db.Prepare("insert " + tbname_econmeeting + " set gid=?,currency=?,title=?,date=?,time=?,time_hidden=?,detail=?,insert_date=?")
			checkErr(err1)
			var n_gid int
			n_gid, _ = strconv.Atoi(result.Id)
			str_date := result.Date
			str_time := result.Time
			str_currency := result.Currency
			str_title := result.Title
			str_timehidden := result.Time_hidden
			str_detail := result.Detail
			str_insert_date := result.Insert_date
			_, err := stmt.Exec(n_gid, str_currency, str_title, str_date, str_time, str_timehidden, str_detail, str_insert_date)
			if err == nil {
				fmt.Printf("====第[%d]组数据开始===\n", index)
				fmt.Printf("id:%s\n", result.Id)
				fmt.Printf("时间:%s %s\n", result.Date, result.Time)
				fmt.Printf("地区:%s\n", result.Currency)
				fmt.Printf("标题:%s\n", result.Title)
				fmt.Printf("====第[%d]组数据结束===\n", index)
			} else {
				fmt.Println(err)
			}
		}
		//=====世界各国假期====
		for index, result := range body.Holiday {
			stmt, err1 := db.Prepare("insert " + tbname_econholiday + " set gid=?,currency=?,title=?,date=?,insert_date=?")
			checkErr(err1)
			var n_gid int
			n_gid, _ = strconv.Atoi(result.Id)
			str_date := result.Date
			str_currency := result.Currency
			str_title := result.Title
			str_insert_date := result.Insert_date
			_, err := stmt.Exec(n_gid, str_currency, str_title, str_date, str_insert_date)
			if err == nil {
				fmt.Printf("====第[%d]组数据开始===\n", index)
				fmt.Printf("id:%s\n", result.Id)
				fmt.Printf("时间:%s %s\n", result.Date)
				fmt.Printf("地区:%s\n", result.Currency)
				fmt.Printf("标题:%s\n", result.Title)
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
	//-------------------
	//开始有数据的日期是2015-01-05
	t := time.Date(2015, time.January, 5, 1, 0, 0, 0, time.Local)
	t_end := time.Now().Add(time.Hour * 24 * 3) //往后获取三天的数据
	//t_end := time.Date(2015, time.June, 5, 1, 0, 0, 0, time.Local)
	//查询一行数据 date
	var str_date_max string
	err = db.QueryRow("SELECT date from " + tbname_econdaily + " ORDER BY date DESC LIMIT 0,10").Scan(&str_date_max)
	if err != nil {
		fmt.Println("query fail") //第一次查询失败 没有数据
		t = time.Date(2015, time.January, 5, 1, 0, 0, 0, time.Local)
	} else {
		fmt.Println("最大日期是:", str_date_max) //查询成功
		t, err = time.Parse("2006-01-02", str_date_max)
		if err != nil {
			fmt.Println("解析时间失败")
		}

		if time.Now().Before(t) {
			t = time.Now() //重新下载
		} else {
			t = t.Add(time.Hour * (-24)) //有可能下载数据中断，往前一天重新下载
		}
	}
	//查询多行
	//rows, errquery := db.Query("SELECT * from " + tbname_econdaily + " ORDER BY date DESC LIMIT 0,10")
	//checkErr(errquery)
	//str_colu, _ := rows.Columns()
	//fmt.Println(str_colu)
	//for rows.Next() == true {
	//	var gid int
	//	var date, time, currency, country, event, period, importance, previous, median, ifr_actual, gjpy, eventpy string
	//	err = rows.Scan(&gid, &date, &time, &currency, &country, &event, &period, &importance, &previous, &median, &ifr_actual, &gjpy, &eventpy)
	//	checkErr(err)
	//	fmt.Println(gid, date, time, currency, country, event, period, importance, previous, median, ifr_actual, gjpy, eventpy)
	//}
	for t_end.After(t) {
		year, month, day := t.Date()
		week := t.Weekday()
		t = t.Add(time.Hour * 24)
		if week == 0 || week == 6 {
			continue
		}
		str_date := fmt.Sprintf("%d-%d-%d", year, month, day)
		readdata(str_date, db)
		//fmt.Println(str_date)
	}
}

var timer_index = 0

func timer1func() {
	fmt.Println("===timer1 run=====", timer_index)
	timer_index++
	//-----数据库操作----
	db, err := sql.Open("mysql", dbopen)
	defer db.Close()
	checkErr(err)
	//-------------------
	//查询数据
	t := time.Now()
	year, month, day := t.Date()
	str_date := fmt.Sprintf("%d-%d-%d", year, month, day)
	updata(str_date, db)
	//var str_time, str_time_min string
	//str_tm1 := time.Now().Format("15:04:05")
	//tm1, _ := time.Parse("15:04:05", str_tm1) //当前时间
	//var tm2 time.Time
	//rows, errquery := db.Query("SELECT time from " + tbname_econdaily + " where date=? and ifr_actual=? ORDER BY time ASC", str_date, "")
	//if errquery != nil {
	//	fmt.Println("查询没有现值的最小时间失败!")
	//	return
	//}
	//defer rows.Close()
	//for rows.Next() == true {
	//	rows.Scan(&str_time)
	//	tm2, _ = time.Parse("15:04:05", str_time) //求没有数据的最小时间值 如果2分钟还没有更新就跳过
	//	tm2 = tm2.Add(time.Minute * 3)
	//	if tm2.After(tm1) == true {
	//		str_time_min = str_time
	//		break
	//	}
	//}
	//fmt.Println(str_time_min)
	//tm2, _ = time.Parse("15:04:05", str_time_min) //没有数据的最小值
	//tm2 = tm2.Add(time.Minute * (-1))
	//fmt.Println(tm1)
	//fmt.Println(tm2)
	//if tm1.After(tm2) == true {
	//	updatarow(str_date, str_time_min, db)
	//}
}
func timer2func() {
	fmt.Println("===timer2 run=====")
	//-----数据库操作----
	db, err := sql.Open("mysql", dbopen)
	defer db.Close()
	checkErr(err)
	//-------------------
	t := time.Date(2015, time.January, 5, 1, 0, 0, 0, time.Local)
	t_end := time.Now().Add(time.Hour * 24 * 3) //往后获取5天的数据
	//查询一行数据 date
	var str_date_max string
	fmt.Sprintf("%d-%d")
	var str_sql string = "SELECT date from " + tbname_econdaily + " ORDER BY date DESC LIMIT 0,10"
	err = db.QueryRow(str_sql).Scan(&str_date_max)
	if err != nil {
		fmt.Println("query fail") //第一次查询失败 没有数据
		t = time.Date(2015, time.January, 5, 1, 0, 0, 0, time.Local)
	} else {
		fmt.Println("最大日期是:", str_date_max) //查询成功
		t, err = time.Parse("2006-01-02", str_date_max)
		if err != nil {
			fmt.Println("解析时间失败")
		}

		if time.Now().Before(t) {
			t = time.Now() //重新下载
		} else {
			t = t.Add(time.Hour * (-24)) //有可能下载数据中断，往前一天重新下载
		}
	}
	for t_end.After(t) {
		year, month, day := t.Date()
		week := t.Weekday()
		t = t.Add(time.Hour * 24)
		if week == 0 || week == 6 {
			continue
		}
		str_date := fmt.Sprintf("%d-%d-%d", year, month, day)
		readdata(str_date, db)
		//fmt.Println(str_date)
	}
}
func starttimer() {
	timer1 := time.NewTicker(1 * time.Minute)
	timer2 := time.NewTicker(12 * time.Hour)
	for {
		select {
		case <-timer1.C:
			timer1func()
		case <-timer2.C:
			timer2func()
		}
	}
}
func main() {
	GetData()
	starttimer()
}
