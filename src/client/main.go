package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

var oauth = struct {
	authURL string
}{
	authURL: "http://192.168.100.101:8080/auth/realms/learningApp/protocol/openid-connect/auth",
}

func main() {
	fmt.Println("OK")
	http.HandleFunc("/", home)
	http.HandleFunc("/login", login)
	http.HandleFunc("/authCodeRedirect", authCodeRedirect)
	http.ListenAndServe(":3000", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("template/index.html"))
	t.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", oauth.authURL, nil)
	if err != nil {
		log.Print(err)
		return
	}
	//req.URL.RawQuery = "state=123abc&client_id=billingApp&response_type=code"
	qs := url.Values{}
	qs.Add("state", "123abc")
	qs.Add("client_id", "billingApp")
	qs.Add("response_type", "code")
	qs.Add("redirect_uri", "http://localhost:3000/authCodeRedirect")

	req.URL.RawQuery = qs.Encode()
	http.Redirect(w, r, req.URL.String(), http.StatusFound)
}

func authCodeRedirect(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Request query: %v", r.URL.Query())

	t := template.Must(template.ParseFiles("template/index.html"))
	t.Execute(w, nil)
}
