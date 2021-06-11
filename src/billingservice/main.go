package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
)

func init() {
	log.SetFlags(log.Ldate + log.Ltime + log.Lmicroseconds + log.LUTC)
}

func main() {
	fmt.Println("OK")
	http.HandleFunc("/", enableLog(home))
	http.ListenAndServe(":4000", nil)
}

func enableLog(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {

	handlerName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	return func(w http.ResponseWriter, r *http.Request) {
		log.SetPrefix(handlerName + " ")
		log.Println("---> " + handlerName)
		log.Printf("request: %v", r.RequestURI)
		handler(w, r)
		log.Println("<--- " + handlerName + "\n")
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "foi")
}
