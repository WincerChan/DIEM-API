package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
)

// HTTPFormat xxxx
type HTTPFormat struct {
	Charset string `json:"charset"`
	Text    string `json:"text"`
}

// query
var hito string
var source string
var content string

// Redirect301 old api
func Redirect301(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/hitokoto/v2/", http.StatusMovedPermanently)
}

func handleError(errorType string) {
	switch errorType {
	case "invalidLengthError":
		hito = "length 参数须为数字且大于 5 哦！"
		source = "Tips"
	default:
		hito = "哦~"
		source = "袴田日向"
	}
}

// FetchAsLength if url param has length
func FetchAsLength(length string) {
	lengthInt, err := strconv.Atoi(length)
	if err != nil || lengthInt < 5 {
		handleError("invalidLengthError")
		return
	}
	err1 := db.QueryRow("SELECT hitokoto, source FROM main WHERE LENGTH(hitokoto) < ? ORDER BY RAND() LiMIT 1;", length).Scan(&hito, &source)
	checkErr(err1)
}

// FetchRandomOne to gen a random int
func FetchRandomOne(length string) {
	if length != "" {
		FetchAsLength(length)
		return
	}
	nBig, err := rand.Int(rand.Reader, big.NewInt(HITOKOTOAMOUNT))
	n := nBig.Int64()
	err = db.QueryRow("SELECT hitokoto, source FROM main LIMIT ?, 1;", n).Scan(&hito, &source)
	checkErr(err)
}

// Hitokoto handle function
func Hitokoto(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	setReqHeader(r)
	// get params
	r.ParseForm()
	encode := r.Form.Get("encode")
	length := r.Form.Get("length")
	callback := r.Form.Get("callback")
	charset := r.Form.Get("charset")
	if charset != "gbk" {
		charset = "utf-8"
	}
	// fetch data
	FetchRandomOne(length)
	// hasCallback is return data
	w.Header().Set("Content-Type", FormatMap["text"].Charset+charset)
	// The value that needs to be returned
	content = fmt.Sprintf(FormatMap["text"].Text, hito, source)
	// set content to encode format
	if text, ok := FormatMap[encode]; ok {
		w.Header().Set("Content-Type", text.Charset+charset)
		content = fmt.Sprintf(text.Text, hito, source)
	}
	// if url params have callback then will ignore encode
	if callback != "" {
		w.Header().Set("Content-Type", "text/javascript; charset="+charset)
		content = fmt.Sprintf("%s({\"hitokoto\": \"%s\", \"source\": \"%s\"})", callback, hito, source)
	}
	// output content
	fmt.Fprint(w, content)

}
