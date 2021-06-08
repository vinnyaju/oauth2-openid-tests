package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

var config = struct {
	authURL             string
	logoutURL           string
	afterLogoutRecirect string
	clientID            string
	authCodeCallback    string
}{
	authURL:             "http://192.168.100.101:8080/auth/realms/learningApp/protocol/openid-connect/auth",
	logoutURL:           "http://192.168.100.101:8080/auth/realms/learningApp/protocol/openid-connect/logout",
	afterLogoutRecirect: "http://localhost:3000/",
	clientID:            "billingApp",
	authCodeCallback:    "http://localhost:3000/authCodeRedirect",
}

//Variáveis privadas da aplicação.
type AppVar struct {
	AuthCode     string
	SessionState string
}

var appVar = AppVar{}

func main() {
	fmt.Println("OK")
	http.HandleFunc("/", home)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/authCodeRedirect", authCodeRedirect)
	http.ListenAndServe(":3000", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("template/index.html"))
	t.Execute(w, appVar)
}

func logout(w http.ResponseWriter, r *http.Request) {

	req, err := http.NewRequest("GET", config.logoutURL, nil)
	if err != nil {
		log.Print(err)
		return
	}

	qs := url.Values{}
	qs.Add("redirect_uri", config.afterLogoutRecirect)

	req.URL.RawQuery = qs.Encode()
	appVar = AppVar{}
	http.Redirect(w, r, req.URL.String(), http.StatusFound)
}

func login(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", config.authURL, nil)
	if err != nil {
		log.Print(err)
		return
	}
	//req.URL.RawQuery = "state=123abc&client_id=billingApp&response_type=code"
	qs := url.Values{}
	qs.Add("state", "123abc")
	qs.Add("client_id", config.clientID)
	qs.Add("response_type", "code")
	qs.Add("redirect_uri", config.authCodeCallback)

	req.URL.RawQuery = qs.Encode()
	http.Redirect(w, r, req.URL.String(), http.StatusFound)
}

func authCodeRedirect(w http.ResponseWriter, r *http.Request) {
	appVar.AuthCode = r.URL.Query().Get("code")
	appVar.SessionState = r.URL.Query().Get("session_state")
	r.URL.RawQuery = ""
	fmt.Printf("Request query: %+v\n", appVar)
	http.Redirect(w, r, "/", http.StatusFound)
}
