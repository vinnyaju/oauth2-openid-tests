package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
)

func init() {
	log.SetFlags(log.Ldate + log.Ltime + log.Lmicroseconds + log.LUTC)
}

type Billing struct {
	Services []string `json:"services"`
}

func main() {
	fmt.Println("OK")
	http.HandleFunc("/", enableLog(home))
	http.HandleFunc("/billing/v1/services", enableLog(services))
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

func services(w http.ResponseWriter, r *http.Request) {
	s := Billing{
		Services: []string{
			"eletricidade",
			"telefonia",
			"internet",
			"água",
		},
	}
	encoder := json.NewEncoder(w)
	w.Header().Add("Content-Type", "application/json")
	encoder.Encode(s)
}
