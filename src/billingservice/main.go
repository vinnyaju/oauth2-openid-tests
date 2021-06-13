package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strings"

	"learn.oauth.billingservice/model"
)

func init() {
	log.SetFlags(log.Ldate + log.Ltime + log.Lmicroseconds + log.LUTC)
}

type Billing struct {
	Services []string `json:"services"`
}

type BillingError struct {
	Error string `json:"error"`
}

var config = struct {
	tokenValidationURI       string
	tokenValidationClientID  string
	tokenValidationClientPWD string
}{
	tokenValidationURI:       "http://192.168.100.101:8080/auth/realms/learningApp/protocol/openid-connect/token/introspect",
	tokenValidationClientID:  "tokenChecker",
	tokenValidationClientPWD: "1bd2dc08-0113-4d23-8b7e-7dcd4ca01965",
}

func main() {
	fmt.Println("OK")
	http.HandleFunc("/", enableLog(home))
	http.HandleFunc("/billing/v1/services", enableLog(utilities))
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

func utilities(w http.ResponseWriter, r *http.Request) {

	token, err := getToken(r)
	if err != nil {
		log.Println(err)
		s := &BillingError{Error: err.Error()}
		encoder := json.NewEncoder(w)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(s)
		return
	}

	log.Println("Token Recebido: ", token)

	//Validar o token
	if !validateToken(token) {
		s := &BillingError{Error: "Token inválido."}
		encoder := json.NewEncoder(w)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(s)
		return
	}

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

func getToken(r *http.Request) (string, error) {
	//Checando o cabeçalho
	token := r.Header.Get("Authorization")

	if token != "" {
		auths := strings.Split(token, " ")
		if len(auths) != 2 {
			return "", fmt.Errorf("invalid Authorization header format")
		}
		return auths[1], nil
	}

	//Checando o corpo do formulário
	token = r.FormValue("access_token")

	if token != "" {
		return token, nil
	}

	//Checando a query string
	token = r.URL.Query().Get("access_token")

	if token != "" {
		return token, nil
	}

	//Retorna erro se não tiver token em lugar nenhum
	return token, fmt.Errorf("access token não informado")
}

func validateToken(token string) bool {

	//Request
	form := url.Values{}
	form.Add("token", token)
	form.Add("token_type_hint", "access_token")

	req, err := http.NewRequest("POST", config.tokenValidationURI, strings.NewReader(form.Encode()))
	req.SetBasicAuth(config.tokenValidationClientID, config.tokenValidationClientPWD)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		log.Print(err)
		return false
	}

	//Client
	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Print(err)
		return false
	}

	//Process Response

	//Check for 200 httpCode return
	if res.StatusCode != 200 {
		log.Println("Status Code returned: ", req.Response.StatusCode)
		log.Println("Status returned: ", req.Response.Status)
		return false
	}

	byteBody, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Print(err)
		return false
	}

	introspect := &model.TokenIntrospect{}
	err = json.Unmarshal(byteBody, introspect)

	if err != nil {
		log.Println(err)
		return false
	}

	return introspect.Active
}
