package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func testTimer1() {
	go func() {
		fmt.Println("test timer1")
	}()

}

func testTimer2() {
	//go func() {
	fmt.Println("test timer2")
	//}()
}

func timer1() {
	timer1 := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timer1.C:
			testTimer1()
		}
	}
}

func timer2() {
	timer2 := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-timer2.C:
			testTimer2()
		}
	}
}
func sayhello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello client! my name is hlm")
}
func main() {
	//go timer1()
	//log.Fatalln("listenandserve:", "123")
	fmt.Println("listenandserve:", "1111")
	http.HandleFunc("/", sayhello)
	err := http.ListenAndServe(":4000", nil)
	//log.Fatalln("listenandserve:", "123")
	timer2()
	//log.Fatalln("listenandserve:", "123")
	if err != nil {
		fmt.Println("listenandserve:", err)
		log.Fatal("listenandserve:", err)
	}
}
