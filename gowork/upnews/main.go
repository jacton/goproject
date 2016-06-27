// test1 project main.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	//"strconv"
	"strings"
	"time"
)

type NewsInfo struct {
	NewsId       int    `json:"NewsId"`
	ChannelId    string `json:"ChannelId"`
	Title        string `json:"Title"`
	Author       string `json:"Author"`
	NewsAbstract string `json:"NewsAbstract"`
	ImgUrl       string `json:"ImgUrl"`
	FileUrl      string `json:"FileUrl"`
	CreateTime   string `json:"conDate"`
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
	timer1 := time.NewTicker(3600 * time.Second)
	for {
		select {
		case <-timer1.C:
			timerfunc()
		}
	}
}

const db_str = "root:123456@/wezone?charset=utf8"

const tb_upnews = "wz_upnews"

func downloadfile(fileurl string, filePath string) {
	resp, err1 := http.Get(fileurl)
	if err1 != nil {
		fmt.Println("get file url fail:", err1)
		return
	}
	defer resp.Body.Close()
	CreateFileDirectory(filePath)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("file create fail:", err)
		return
	}
	defer file.Close()
	_, err2 := io.Copy(file, resp.Body)
	if err2 != nil {
		fmt.Println("write file fail:", err2)
	} else {
		fmt.Println("down file[" + fileurl + "]success")
	}
}
func readdata() {
	fmt.Println("==start get url==")
	resp, errHttp := http.Get("http://app.upchinafund.com/sasweb/xysidkdydnhensydn_cdhds.dyshg/dsfyewlrndsfpoidsfewlkdsnf.cxgdsf_hdsfnew_gz/AjaxFRInews/GetFinancialIntelligenceByNumCache.cspx?channelId=041001&rowIndex=0&selectNumber=100&jsoncallback=jQuery180011304811830632389_1448590242312&_=1448590242338")
	if errHttp != nil {
		fmt.Println("get web fail")
		return
	}
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	strBody := string(data)
	strBody = strings.TrimLeft(strBody, "jQuery180011304811830632389_1448590242312(")
	strBody = strings.TrimRight(strBody, ")")
	//var body Body_info
	//-----数据库操作----
	db, err := sql.Open("mysql", db_str)
	defer db.Close()
	checkErr(err)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " + tb_upnews + " ( " +
		"`id` int(11) NOT NULL AUTO_INCREMENT," +
		"`newsid` int(11) DEFAULT 0," +
		"`channelid` varchar(255) DEFAULT NULL," +
		"`title` varchar(255) DEFAULT NULL," +
		"`author` varchar(255) DEFAULT NULL," +
		"`newsabstract` text," +
		"`filebaseurl` varchar(255) DEFAULT NULL," +
		"`imgurl` varchar(255) DEFAULT NULL," +
		"`fileurl` varchar(255) DEFAULT NULL," +
		"`createtime` datetime DEFAULT NULL," +
		"PRIMARY KEY (`id`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8;")
	checkErr(err)
	var tid int
	err = db.QueryRow("SELECT MAX(newsid) FROM " + tb_upnews).Scan(&tid)
	checkErr(err)
	fmt.Println("max id:", tid)
	//return
	//-----------------
	var body []NewsInfo
	if err := json.Unmarshal([]byte(strBody), &body); err == nil {
		count := 0
		for index := len(body) - 1; index >= 0; index-- {
			result := body[index]
			if result.NewsId <= tid {
				continue
			}
			count++
			fmt.Printf("====第[%d]组数据开始===\n", index)
			fmt.Printf("newsid:%d\n", result.NewsId)
			fmt.Printf("时间:%s\n", result.CreateTime)
			fmt.Printf("标题:%s\n", result.Title)
			fmt.Printf("内容:%s\n", result.NewsAbstract)
			filebaseurl := "http://img.upchinafund.com"
			fileurl := result.FileUrl
			imgurl := result.ImgUrl
			fmt.Printf("====第[%d]组数据结束===\n", index)
			stmt, err1 := db.Prepare("insert " + tb_upnews + " set newsid=?,channelid=?,title=?,author=?,newsabstract=?,filebaseurl=?,imgurl=?,fileurl=?,createtime=?")
			checkErr(err1)
			_, err = stmt.Exec(result.NewsId, result.ChannelId, result.Title, result.Author, result.NewsAbstract, filebaseurl, imgurl, fileurl, result.CreateTime)
			checkErr(err)
			//下载图片和文件
			if imgurl != "" {
				imgpath := GetCurrPath() + imgurl
				downloadfile(filebaseurl+imgurl, imgpath)
			}
			if fileurl != "" {
				pdfpath := GetCurrPath() + fileurl
				downloadfile(filebaseurl+fileurl, pdfpath)
			}
			//==========
		}
		fmt.Println("get ", count, " data")
	} else {
		fmt.Println("解析json失败:", err)
	}
}
func GetCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	splitstring := strings.Split(path, "\\")
	size := len(splitstring)
	splitstring = strings.Split(path, splitstring[size-1])
	ret := strings.Replace(splitstring[0], "\\", "/", size-1)
	return ret
}
func CreateFileDirectory(filepath string) {
	splitstring := strings.Split(filepath, "/")
	size := len(splitstring)
	splitstring = strings.Split(filepath, splitstring[size-1])
	os.MkdirAll(splitstring[0], 0777)
}
func main() {
	readdata()
	starttimer()
}
