// web project main.go
package main

import (
	"fmt"
	"html/template"
	"io"
	//"io/ioutil"
	"log"
	//"mime/multipart"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func sayhello(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	//fmt.Println(r.Form)
	//fmt.Println("path", r.URL.Path)
	//fmt.Println("scheme", r.URL.Scheme)
	//fmt.Println(r.Form["url_log"])
	//for k, v := range r.Form {
	//	fmt.Println("key:", k)
	//	fmt.Println("val:", strings.Join(v, ","))
	//}
	//fmt.Fprintf(w, "hello client! my name is hlm")
	fmt.Println(r.URL.Path)
	if ok, _ := regexp.MatchString("/img/", r.URL.String()); ok {
		http.StripPrefix("/img", http.FileServer(http.Dir("./img/"))).ServeHTTP(w, r)
	} else {
		fmt.Fprintf(w, "hello client! my name is hlm")
	}
}
func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ","))
	}
	fmt.Println("----login------")
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.tpl")
		t.Execute(w, nil)
		fmt.Println("get execute")
	} else {
		fmt.Println("rcv post msg")
		fmt.Println("username:", r.Form["username"])
		fmt.Println("pwd:", r.Form["userpwd"])
		fmt.Fprintf(w, "longin sucess")
	}
}
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-----upload-----")
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		fmt.Println("rcv get msg")
		t, _ := template.ParseFiles("upfile.tpl")
		t.Execute(w, nil)
	} else {
		fmt.Println("rcv post msg")
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println("解析失败 formfile")
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./img/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("写文件失败")
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}
func img(w http.ResponseWriter, r *http.Request) {
	fmt.Println("----img------")
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		http.FileServer(http.Dir("./img/")).ServeHTTP(w, r)
		//http.StripPrefix("/img", http.FileServer(http.Dir("./img/"))).ServeHTTP(w, r)
	}
}
func mysql(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ","))
	}
	fmt.Println("-----mysql-----")
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.tpl")
		t.Execute(w, nil)
		fmt.Println("get execute")
	} else {
		fmt.Println("===0====")
		db, err := sql.Open("mysql", "root:123456@/test?charset=utf8")
		fmt.Println("===1====")
		checkErr(err)
		fmt.Println(r.Form["username"])
		stmt, err := db.Prepare("insert student set id=?,name=?,age=?")
		fmt.Println("===2====")
		checkErr(err)
		name := r.FormValue("username")
		res, err := stmt.Exec("2", name, "18")
		id, err := res.LastInsertId()
		fmt.Println("===3====")
		checkErr(err)
		fmt.Println(id)
		//fmt.Println("rcv post msg")
		//fmt.Println("username:", r.Form["username"])
		//fmt.Println("pwd:", r.Form["userpwd"])
		fmt.Fprintf(w, "insert my sql sucess")
	}
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	http.HandleFunc("/", sayhello)
	http.HandleFunc("/login", login)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/mysql", mysql)
	http.HandleFunc("/img", img)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("listenandserve:", err)
		log.Fatal("listenandserve:", err)
	} else {
		fmt.Println("listen on 9090:")
	}

}
