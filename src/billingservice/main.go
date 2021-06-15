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
	Utilities []string `json:"utilities"`
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
	http.HandleFunc("/billing/v1/utilities", enableLog(utilities))
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
		buildErrorMessage(w, err.Error())
		return
	}

	//log.Println("Token Recebido: ", token)

	//Validar o token
	tokenValid, tokenClaim := validateTokenAndExtractClaim(token)
	if !tokenValid {
		buildErrorMessage(w, "Token inválido.")
		return
	}

	//Validar se há o escopo necessário à execução deste serviço
	//log.Println("Scopes: ", tokenClaim.Scope)
	if !strings.Contains(tokenClaim.Scope, "getUtilitiesService") {
		buildErrorMessage(w, "O escopo necessário à execução deste serviço não está presente no token informado.")
		return
	}

	s := Billing{
		Utilities: []string{
			"eletricidade",
			"telefonia",
			"internet",
			"água",
		},
	}
	encoder := json.NewEncoder(w)
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	encoder.Encode(s)
}

func getToken(r *http.Request) (string, error) {
	//Checando o cabeçalho
	token := r.Header.Get("Authorization")

	if token != "" {
		auths := strings.Split(token, " ")
		if len(auths) != 2 {
			return "", fmt.Errorf("Formato do cabeçalho Authorization inválido")
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
	return token, fmt.Errorf("Access token não informado")
}

func validateTokenAndExtractClaim(token string) (bool, *model.TokenIntrospect) {

	//Request
	form := url.Values{}
	form.Add("token", token)
	form.Add("token_type_hint", "access_token")

	req, err := http.NewRequest("POST", config.tokenValidationURI, strings.NewReader(form.Encode()))
	req.SetBasicAuth(config.tokenValidationClientID, config.tokenValidationClientPWD)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		log.Print(err)
		return false, nil
	}

	//Client
	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Print(err)
		return false, nil
	}

	//Process Response

	//Check for 200 httpCode return
	if res.StatusCode != 200 {
		log.Println("Status Code returned: ", req.Response.StatusCode)
		log.Println("Status returned: ", req.Response.Status)
		return false, nil
	}

	byteBody, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Print(err)
		return false, nil
	}

	introspect := &model.TokenIntrospect{}
	err = json.Unmarshal(byteBody, introspect)

	if err != nil {
		log.Println(err)
		return false, nil
	}

	return introspect.Active, introspect
}

func buildErrorMessage(w http.ResponseWriter, message string) {
	s := &BillingError{Error: message}
	encoder := json.NewEncoder(w)
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusBadRequest)
	encoder.Encode(s)
}
