package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strings"

	"learn.oauth.client/model"
)

var config = struct {
	authURL             string
	logoutURL           string
	afterLogoutRecirect string
	clientID            string
	clientPWD           string
	authCodeCallback    string
	tokenIssuerURL      string
}{
	authURL:             "http://192.168.100.101:8080/auth/realms/learningApp/protocol/openid-connect/auth",
	logoutURL:           "http://192.168.100.101:8080/auth/realms/learningApp/protocol/openid-connect/logout",
	tokenIssuerURL:      "http://192.168.100.101:8080/auth/realms/learningApp/protocol/openid-connect/token",
	afterLogoutRecirect: "http://localhost:3000/",
	clientID:            "billingApp",
	clientPWD:           "3bd73711-a702-4494-82ea-d280ea5a855c",
	authCodeCallback:    "http://localhost:3000/authCodeRedirect",
}

//Variáveis privadas da aplicação.
type AppVar struct {
	AuthCode     string
	SessionState string
	AccessToken  string
	RefreshToken string
	Scope        string
}

var appVar = AppVar{}

func init() {
	log.SetFlags(log.Ldate + log.Ltime + log.Lmicroseconds + log.LUTC)
}

func main() {
	fmt.Println("OK")
	http.HandleFunc("/", enableLog(home))
	http.HandleFunc("/login", enableLog(login))
	http.HandleFunc("/exchangeToken", enableLog(exchangeToken))
	http.HandleFunc("/logout", enableLog(logout))
	http.HandleFunc("/authCodeRedirect", enableLog(authCodeRedirect))
	http.ListenAndServe(":3000", nil)
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

func exchangeToken(w http.ResponseWriter, r *http.Request) {
	//Request
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", appVar.AuthCode)
	form.Add("redirect_uri", config.authCodeCallback)
	form.Add("client_id", config.clientID)
	req, err := http.NewRequest("POST", config.tokenIssuerURL, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		log.Print(err)
		return
	}

	req.SetBasicAuth(config.clientID, config.clientPWD)
	//Client
	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Print(err)
		return
	}

	//Process response
	byteBody, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Print(err)
		return
	}

	accessTokenResponse := &model.AccessTokenResponse{}

	json.Unmarshal(byteBody, accessTokenResponse)

	appVar.AccessToken = accessTokenResponse.AccessToken
	appVar.RefreshToken = accessTokenResponse.RefreshToken
	appVar.Scope = accessTokenResponse.Scope
	log.Printf("AccessToken: %v", appVar.AccessToken)
	log.Printf("RefreshToken: %v", appVar.RefreshToken)
	log.Printf("Scope: %v", appVar.Scope)

	http.Redirect(w, r, "/", http.StatusFound)
	// t := template.Must(template.ParseFiles("template/index.html"))
	// t.Execute(w, appVar)
}
