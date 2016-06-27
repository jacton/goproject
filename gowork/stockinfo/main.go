// test1 project main.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Body_info struct {
	Rc      int        `json:"rc"`
	Me      string     `json:"me"`
	Results []NewsInfo `json:"LivesList"`
}
type NewsInfo struct {
	//Id         string `json:"id"`
	Newsid string `json:"newsid"`
	Url_w  string `json:"url_w"`
	//Url_m     string `json:"url_m"`
	Title string `json:"title"`
	//Simtitle  string `json:"simtitle"`
	Digest string `json:"digest"`
	//Simdigest string `json:"simdigest"`
	//Image string `json:"image"`
	//Titlestyle string `json:"titlestyle"`
	//Url_pdf    string `json:"url_pdf"`
	//Type1      string `json:"type"`
	//Simtype    string `json:"simtype"`
	//Simtype_zh string `json:"simtype_zh"`
	//topic           string `json:"topic"`
	//simspecial      string `json:"simspecial"`
	//simtopic        string `json:"simtopic"`
	//column          string `json:"column"`
	//Editor_id       string `json:"editor_id"`
	//Editor_name     string `json:"editor_name"`
	//Lasteditor_name string `json:"lasteditor_name"`
	Showtime string `json:"showtime"`
	//Ordertime       string `json:"ordertime"`
	//commentnum      string `json:"commentnum"`
	//newstype string `json:"newstype"`
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
}

func timerfunc() {
	fmt.Println("===timer run=====")
	readdata()
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

const db_str = "root:123456@/news?charset=utf8"

const tb_stocknews = "stocknews"

func readdata() {
	fmt.Println("==start read data==")
	resp, errhttp := http.Get("http://newsapi.eastmoney.com/kuaixun/v1/getlist_zhiboall_ajaxResult_70_1_.html")
	if errhttp != nil {
		fmt.Println("get web fail")
		return
	}
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	//-----数据库操作----
	db, err := sql.Open("mysql", db_str)
	defer db.Close()
	checkErr(err)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " + tb_stocknews + "(" +
		"`id` int(11) NOT NULL AUTO_INCREMENT," +
		"`gid` bigint(20) DEFAULT NULL," +
		"`time` datetime DEFAULT NULL," +
		"`title` text CHARACTER SET utf8," +
		"`content` text CHARACTER SET utf8," +
		"`url` varchar(255) CHARACTER SET utf8," +
		"PRIMARY KEY (`id`)" +
		") ENGINE=InnoDB AUTO_INCREMENT=52 DEFAULT CHARSET=utf8;")
	checkErr(err)
	var tid int64
	err = db.QueryRow("SELECT MAX(gid) FROM " + tb_stocknews).Scan(&tid)
	checkErr(err)
	fmt.Println("max id:", tid)
	//return
	//-----------------
	var str_data string
	str_data = string(data)
	str_data = strings.TrimLeft(str_data, "var ajaxResult=")
	var body Body_info
	if errjson := json.Unmarshal([]byte(str_data), &body); errjson == nil {
		count := 0
		for index := len(body.Results) - 1; index >= 0; index-- {
			result := body.Results[index]
			id, _ := strconv.ParseInt(result.Newsid, 10, 64)
			if id <= tid {
				continue
			}
			count++
			fmt.Printf("====第[%d]组数据开始===\n", index)
			fmt.Printf("id:%s\n", result.Newsid)
			fmt.Printf("时间:%s\n", result.Showtime)
			fmt.Printf("标题:%s\n", result.Title)
			fmt.Printf("内容:%s\n", result.Digest)
			fmt.Printf("链接:%s\n", result.Url_w)
			fmt.Printf("====第[%d]组数据结束===\n", index)
			stmt, err1 := db.Prepare("insert " + tb_stocknews + " set gid=?,time=?,title=?,content=?,url=?")
			checkErr(err1)
			str_time := result.Showtime
			str_title := result.Title
			str_content := result.Digest
			str_url := result.Url_w
			_, err = stmt.Exec(id, str_time, str_title, str_content, str_url)
			checkErr(err)
		}
		fmt.Println("get ", count, " data")
	} else {
		fmt.Println("解析json失败:", errjson)
	}
}
func main() {
	readdata()
	starttimer()
}
