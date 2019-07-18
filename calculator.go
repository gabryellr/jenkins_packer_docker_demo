package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

var historic []string

type CalJson struct {
	Result float64
}

type ErrorJson struct {
	Error string
}

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/calc/{operation}/{a}/{b}", httpHandle)
	mux.HandleFunc("/calc/history", historicMath)
	http.Handle("/", mux)
	http.ListenAndServe(":8080", nil)
}

func httpHandle(httpResponse http.ResponseWriter, httpRequest *http.Request) {
	httpResponse.Header().Set("Content-Type", "application/json")
	urlPart := strings.Split(httpRequest.URL.Path, "/")
	firstNumber, err1 := strconv.ParseFloat(urlPart[3], 10)
	secondNumber, err2 := strconv.ParseFloat(urlPart[4], 10)

	if err1 == nil && err2 == nil {
		switch urlPart[2] {
		case "sum":
			historic = append(historic, urlPart[3]+"+"+urlPart[4])
			js, err := json.Marshal(CalJson{sum(firstNumber, secondNumber)})
			if err != nil {
				http.Error(httpResponse, err.Error(), http.StatusInternalServerError)
				return
			}
			httpResponse.Write(js)
		case "sub":
			historic = append(historic, urlPart[3]+"-"+urlPart[4])
			js, err := json.Marshal(CalJson{sub(firstNumber, secondNumber)})
			if err != nil {
				http.Error(httpResponse, err.Error(), http.StatusInternalServerError)
				return
			}
			httpResponse.Write(js)
		case "mul":
			historic = append(historic, urlPart[3]+"*"+urlPart[4])
			js, err := json.Marshal(CalJson{mul(firstNumber, secondNumber)})
			if err != nil {
				http.Error(httpResponse, err.Error(), http.StatusInternalServerError)
				return
			}
			httpResponse.Write(js)
		case "div":
			historic = append(historic, urlPart[3]+"/"+urlPart[4])
			js, err := json.Marshal(CalJson{div(firstNumber, secondNumber)})
			if err != nil {
				http.Error(httpResponse, err.Error(), http.StatusInternalServerError)
				return
			}
			httpResponse.Write(js)
		default:
			js, err := json.Marshal(ErrorJson{"invalid operator"})
			if err != nil {
				http.Error(httpResponse, err.Error(), http.StatusInternalServerError)
				return
			}
			jsonError := []byte(js)
			httpResponse.WriteHeader(http.StatusBadRequest)
			httpResponse.Write([]byte(jsonError))
		}
	} else {
		js, err := json.Marshal(ErrorJson{"Type the correct path"})
		if err != nil {
			http.Error(httpResponse, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonError := []byte(js)
		httpResponse.WriteHeader(http.StatusBadRequest)
		httpResponse.Write([]byte(jsonError))
	}
}

func historicMath(httpResponse http.ResponseWriter, req *http.Request) {
	urlsJson, _ := json.Marshal(historic)
	fmt.Fprint(httpResponse, string(urlsJson))
}

func sum(firstNumber, secondNumber float64) float64 {
	return firstNumber + secondNumber
}

func sub(firstNumber, secondNumber float64) float64 {
	return firstNumber - secondNumber
}

func mul(firstNumber, secondNumber float64) float64 {
	return firstNumber * secondNumber
}

func div(firstNumber, secondNumber float64) float64 {
	return firstNumber / secondNumber
}