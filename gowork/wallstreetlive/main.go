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
	"time"
)

type Paginator struct {
	Total    int    `json:"total"`
	Previous string `json:"previous"`
	Next     string `json:"next"`
	Last     string `json:"last"`
}
type Text struct {
	ContentExtra    string `json:"contentExtra"`
	ContentFollowup string `json:"contentFollowup"`
	ContentAnalysis string `json:"contentAnalysis"`
}
type Body_info struct {
	//Paginator Paginator  `json:"paginator"`
	Results []NewsInfo `json:"results"`
}
type NewsInfo struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	//Type1      string `json:"type"`
	//CodeType   string `json:"codeType"`
	//Importance string `json:"importance"`
	CreatedAt string `json:"createdAt"`
	//UpdatedAt  string `json:"updatedAt"`
	//ImageCount    string `json:"imageCount"`
	//Image         string `json:"image"`
	//VideoCount    string `json:"videoCount"`
	//Video         string `json:"video"`
	//ViewCount     string `json:"viewCount"`
	//ShareCount    string `json:"shareCount"`
	CommentStatus string `json:"commentStatus"`
	ContentHtml   string `json:"contentHtml"`
	//Data          string `json:"data"`
	//SourceName    string `json:"sourceName"`
	//SourceUrl     string `json:"sourceUrl"`
	//UserId        string `json:"userId"`
	//CategorySet   string `json:"categorySet"`
	//HasMore       string `json:"hasMore"`
	//ChannelSet    string `json:"channelSet"`
	//Text          Text   `json:"text"`
	Node_color string `json:"node_color"`
	//Node_format string `json:"node_format"`
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

const tb_fastnews = "fastnews"

func readdata() {
	fmt.Println("==start read data==")
	resp, errhttp := http.Get("http://api.wallstreetcn.com/v2/livenews")
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
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " + tb_fastnews + "(" +
		"`id` int(11) NOT NULL AUTO_INCREMENT," +
		"`gid` int(11) DEFAULT NULL," +
		"`time` datetime DEFAULT NULL," +
		"`title` text CHARACTER SET utf8," +
		"`content` text CHARACTER SET utf8," +
		"`color` varchar(255) DEFAULT NULL," +
		"PRIMARY KEY (`id`)" +
		") ENGINE=InnoDB AUTO_INCREMENT=52 DEFAULT CHARSET=utf8;")
	checkErr(err)
	var tid int
	err = db.QueryRow("SELECT MAX(gid) FROM " + tb_fastnews).Scan(&tid)
	checkErr(err)
	fmt.Println("max id:", tid)
	//-----------------
	var body Body_info
	if err := json.Unmarshal([]byte(data), &body); err == nil {
		count := 0
		for index := len(body.Results) - 1; index >= 0; index-- {
			result := body.Results[index]
			id, _ := strconv.Atoi(result.Id)
			if id <= tid {
				continue
			}
			count++
			fmt.Printf("====第[%d]组数据开始===\n", index)
			fmt.Printf("id:%s\n", result.Id)
			t, _ := strconv.ParseInt(result.CreatedAt, 10, 64)
			//2006-01-02 15:04:05
			str_time := time.Unix(t, 0).Format("2006-01-02 15:04:05")
			fmt.Printf("时间:%s\n", str_time)
			fmt.Printf("标题:%s\n", result.Title)
			fmt.Printf("内容:%s\n", result.ContentHtml)
			fmt.Printf("颜色:%s\n", result.Node_color)
			fmt.Printf("====第[%d]组数据结束===\n", index)
			stmt, errdb := db.Prepare("insert " + tb_fastnews + " set gid=?,time=?,title=?,content=?,color=?")
			checkErr(errdb)
			_, err = stmt.Exec(id, str_time, result.Title, result.ContentHtml, result.Node_color)
			checkErr(err)
		}
		fmt.Println("get ", count, " data")
	} else {
		fmt.Println("解析json失败:", err)
	}
}
func main() {
	readdata()
	starttimer()
}
