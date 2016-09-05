package models

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	//"time"
	"github.com/astaxie/beego"
)

const (
	ret_success = iota
	ret_fail
	ret_db_error
)

var (
	dbnews *sql.DB
)

var db_str string

var tb_fastnews string

var tb_jinshinews string

func init() {
	db_str = beego.AppConfig.String("db_username") + ":" + beego.AppConfig.String("db_password") + "@/" + beego.AppConfig.String("db_name") + "?charset=utf8"
	tb_fastnews = beego.AppConfig.String("tb_fastnews")
	tb_jinshinews = beego.AppConfig.String("tb_jinshinews")
	dbnews, _ = sql.Open("mysql", db_str)
	fmt.Println("init models")
}

type News struct {
	Articid  string
	Title    string
	Content  string
	Datetime string
	Color    string
}

type NewsInfo struct {
	Type        string
	Name        string
	Articid     string
	Title       string
	Abstract    string
	Content     string
	Contenthtml string
	Detailurl   string
	Datetime    string
	Autor       string
	Imgname     string
	Imgurl      string
}

func GetNews(id string, limit string, direction string) map[string]interface{} {
	jsdata := make(map[string]interface{})
	fmt.Println("GetNews id is", id)
	if dbnews == nil {
		fmt.Println("GetNews connect datebase")
		dbnews, _ = sql.Open("mysql", db_str)
	}
	if dbnews == nil {
		jsdata["retcode"] = ret_db_error
		jsdata["retmsg"] = "db connect error"
		return jsdata
	}
	var rows *sql.Rows
	var err error
	if limit == "" {
		limit = "30"
	}
	if direction == "" {
		direction = "0"
	}
	if id == "" || id == "0" {
		rows, err = dbnews.Query("SELECT * from "+tb_fastnews+" ORDER BY gid DESC LIMIT ?", limit)
	} else {
		if direction == "0" {
			rows, err = dbnews.Query("SELECT * from (SELECT * from "+tb_fastnews+" WHERE  gid >?  LIMIT ?) as t ORDER BY gid DESC ", id, limit)
		} else {
			rows, err = dbnews.Query("SELECT * from "+tb_fastnews+" WHERE  gid <?  ORDER BY gid DESC LIMIT ?", id, limit)
		}

	}
	if err == nil {
		var id, gid int
		var str_title, str_content, str_time, str_color string
		var arrNews []News
		for rows.Next() == true {
			rows.Scan(&id, &gid, &str_time, &str_title, &str_content, &str_color)
			var tmpnews News
			tmpnews.Title = str_title
			tmpnews.Datetime = str_time
			tmpnews.Content = str_content
			tmpnews.Color = str_color
			tmpnews.Articid = strconv.Itoa(gid)
			arrNews = append(arrNews, tmpnews)
		}
		jsdata["retcode"] = ret_success
		jsdata["retmsg"] = "success"
		jsdata["data"] = arrNews
	} else {
		jsdata["retcode"] = ret_fail
		jsdata["retmsg"] = "fail"
	}
	return jsdata
}
func GetNewsInfo(strtype string, strdatetime string, limit string, direction string) map[string]interface{} {
	jsdata := make(map[string]interface{})
	fmt.Println("GetNewsinfo type is ", strtype, " time is ", strdatetime, " direction is ", direction)
	if dbnews == nil {
		fmt.Println("GetNewsinfo connect datebase")
		dbnews, _ = sql.Open("mysql", db_str)
	}
	if dbnews == nil {
		jsdata["retcode"] = ret_db_error
		jsdata["retmsg"] = "db connect error"
		return jsdata
	}
	var rows *sql.Rows
	var err error
	if limit == "" {
		limit = "30"
	}
	if direction == "" {
		direction = "0"
	}
	if strtype == "" && strdatetime == "" {
		rows, err = dbnews.Query("SELECT type,`name`,articid,title,abstract,content,contenthtml,time,autor,imgname,imgurl FROM "+tb_jinshinews+" ORDER BY time DESC LIMIT ?", limit)
	} else if strtype != "" && strdatetime == "" {
		rows, err = dbnews.Query("SELECT type,`name`,articid,title,abstract,content,contenthtml,time,autor,imgname,imgurl FROM "+tb_jinshinews+" WHERE type = ? ORDER BY time DESC LIMIT ?", strtype, limit)
	} else if strtype == "" && strdatetime != "" {
		rows, err = dbnews.Query("SELECT type,`name`,articid,title,abstract,content,contenthtml,time,autor,imgname,imgurl FROM "+tb_jinshinews+" WHERE time > ? ORDER BY time DESC LIMIT ?", strdatetime, limit)
	} else if strtype != "" && strdatetime != "" && direction == "0" {
		rows, err = dbnews.Query("SELECT type,`name`,articid,title,abstract,content,contenthtml,time,autor,imgname,imgurl FROM "+tb_jinshinews+" WHERE time > ? and type=? ORDER BY time DESC LIMIT ?", strdatetime, strtype, limit)
	} else if strtype != "" && strdatetime != "" && direction != "0" {
		rows, err = dbnews.Query("SELECT type,`name`,articid,title,abstract,content,contenthtml,time,autor,imgname,imgurl FROM "+tb_jinshinews+" WHERE time < ? and type=? ORDER BY time DESC LIMIT ?", strdatetime, strtype, limit)
	}
	if err == nil {
		var str_type, str_name, str_articid, str_title, str_abstract, str_content, str_contenthtml, str_time, str_autor, str_imgname, str_imgurl string
		var arrNewsInfo []NewsInfo
		for rows.Next() == true {
			rows.Scan(&str_type, &str_name, &str_articid, &str_title, &str_abstract, &str_content, &str_contenthtml, &str_time, &str_autor, &str_imgname, &str_imgurl)
			var tmpnews NewsInfo
			tmpnews.Type = str_type
			tmpnews.Name = str_name
			tmpnews.Articid = str_articid
			tmpnews.Title = str_title
			tmpnews.Abstract = str_abstract
			//tmpnews.Content = str_content
			//tmpnews.Contenthtml = str_contenthtml
			tmpnews.Datetime = str_time
			tmpnews.Autor = str_autor
			tmpnews.Imgname = str_imgname
			tmpnews.Imgurl = str_imgurl
			arrNewsInfo = append(arrNewsInfo, tmpnews)
		}
		jsdata["retcode"] = ret_success
		jsdata["retmsg"] = "success"
		jsdata["data"] = arrNewsInfo
	} else {
		jsdata["retcode"] = ret_fail
		jsdata["retmsg"] = "fail"
	}
	return jsdata
}
func GetNewsInfoById(strid string) map[string]interface{} {
	jsdata := make(map[string]interface{})
	fmt.Println("GetNewsbyid id is", strid)
	if dbnews == nil {
		fmt.Println("GetNewsbyid connect datebase")
		dbnews, _ = sql.Open("mysql", db_str)
	}
	if dbnews == nil {
		jsdata["retcode"] = ret_db_error
		jsdata["retmsg"] = "db connect error"
		return jsdata
	}
	var rows *sql.Rows
	var err error
	if strid == "" {
		err = errors.New("strid is empty")
	}
	rows, err = dbnews.Query("SELECT type,`name`,articid,title,abstract,content,contenthtml,time,autor,imgname,imgurl FROM "+tb_jinshinews+" where articid = ? ORDER BY time DESC LIMIT 2", strid)
	if err == nil {
		var str_type, str_name, str_articid, str_title, str_abstract, str_content, str_contenthtml, str_time, str_autor, str_imgname, str_imgurl string
		var arrNewsInfo []NewsInfo
		for rows.Next() == true {
			rows.Scan(&str_type, &str_name, &str_articid, &str_title, &str_abstract, &str_content, &str_contenthtml, &str_time, &str_autor, &str_imgname, &str_imgurl)
			var tmpnews NewsInfo
			tmpnews.Type = str_type
			tmpnews.Name = str_name
			tmpnews.Articid = str_articid
			tmpnews.Title = str_title
			tmpnews.Abstract = str_abstract
			tmpnews.Content = str_content
			tmpnews.Contenthtml = str_contenthtml
			tmpnews.Datetime = str_time
			tmpnews.Autor = str_autor
			tmpnews.Imgname = str_imgname
			tmpnews.Imgurl = str_imgurl
			tmpnews.Detailurl = beego.AppConfig.String("tb_jinshinews_url") + str_articid
			arrNewsInfo = append(arrNewsInfo, tmpnews)
		}
		jsdata["retcode"] = ret_success
		jsdata["retmsg"] = "success"
		jsdata["data"] = arrNewsInfo
	} else {
		jsdata["retcode"] = ret_fail
		jsdata["retmsg"] = "fail"
	}
	return jsdata
}
