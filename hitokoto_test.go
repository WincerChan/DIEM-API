package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ReqTest func(w http.ResponseWriter, r *http.Request)

func init() {
	initConfig("./config.yaml")
	initHitokotoDB()
	initRedis()
	MakeReturnMap()
}

func TestHitokoto(t *testing.T) {
	urlParamsTest := func(reqFunc ReqTest, url string) {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		checkErr(err)
		rw := httptest.NewRecorder()
		handler := http.HandlerFunc(reqFunc)
		handler.ServeHTTP(rw, req)
		log.Println("code:", rw.Code)
		log.Println("body:", rw.Body.String())
	}
	urlParamsTest(Hitokoto, "/hitokoto/v2/?encode=json")
	urlParamsTest(Hitokoto, "/hitokoto/v2/?callback=fff&length=4")
	urlParamsTest(Hitokoto, "/hitokoto/v2/?length=14")
	urlParamsTest(Redirect301, "/hitokoto/get")
	//
}
