package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

// query
var hitoinfo *Info
var hito string
var source string
var content string
var pipe chan string

type Info struct {
	Source string `json:"source"`
	Hito   string `json:"hitokoto"`
}

func (c Info) Value() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Info) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}

func setLimitHeader(w http.ResponseWriter, r *http.Request) bool {
	// If user do not enable redis-cell limit, just do not check
	if !config.Redis.Enabled {
		return false
	}
	ret := getRemainingNumbers(r)
	w.Header().Set("X-RateLimit-Limit", strconv.FormatInt(ret[1].(int64), 10))
	w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(ret[2].(int64), 10))
	w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(ret[4].(int64), 10))
	if ret[0].(int64) == 1 {
		content = "{\"result\": \"Your IP requests is frequently.\"}"
		w.Write([]byte(content))
		return true
	}
	return false
}

func parseParams(r *http.Request) {
	r.ParseForm()
	pipe <- r.Form.Get("charset")
	pipe <- r.Form.Get("encode")
	pipe <- r.Form.Get("length")
	pipe <- r.Form.Get("callback")
}

func fetchInfo() {
	lenStr := <-pipe
	if lenStr == "" || len(lenStr) > 3 {
		db.QueryRow("SELECT RANDOMFETCH($1);", -1).Scan(&hitoinfo)
	} else {
		length, err := strconv.Atoi(lenStr)
		checkErr(err)
		db.QueryRow("SELECT RANDOMFETCH($1);", length).Scan(&hitoinfo)
	}
}

func setResponse(w http.ResponseWriter) {
	var buffer bytes.Buffer
	charset := <-pipe
	contenttype := ""
	if "gbk" != charset {
		charset = "utf-8"
	}
	switch e := <-pipe; e {
	case "js":
		contenttype = "text/javascript; charset=" + charset
		buffer.WriteString("var hitokoto=\"")
		buffer.WriteString(hitoinfo.Hito)
		buffer.WriteString("——「")
		buffer.WriteString(hitoinfo.Source)
		buffer.WriteString("」\";var dom=document.querySelector('.hitokoto');")
		buffer.WriteString("Array.isArray(dom)?dom[0].innerText=hitokoto:dom.innerText=hitokoto;")
		break
	case "json":
		contenttype = "application/json; charset=" + charset
		buffer.WriteString("{\"hitokoto\": \"")
		buffer.WriteString(hitoinfo.Hito)
		buffer.WriteString("\", \"source\": \"")
		buffer.WriteString(hitoinfo.Source)
		buffer.WriteString("\"}")
		break
	default:
		contenttype = "text/plain; charset=" + charset
		buffer.WriteString(hitoinfo.Hito)
		buffer.WriteString("——「")
		buffer.WriteString(hitoinfo.Source)
		buffer.WriteString("」")
	}
	w.Header().Set("Content-Type", contenttype)
	w.Write(buffer.Bytes())
}

func setCallback(w http.ResponseWriter) {
	json, _ := hitoinfo.Value()
	var buffer bytes.Buffer
	callback := <-pipe
	if "" == callback {
		return
	}
	contenttype := "text/javascript"
	w.Header().Set("Content-Type", contenttype)
	buffer.WriteString(callback)
	buffer.WriteString("(")
	buffer.WriteString(string(json))
	buffer.WriteString(")")
	w.Write(buffer.Bytes())

}

// Hitokoto handle function
func Hitokoto(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	isLimited := setLimitHeader(w, r)
	if isLimited {
		return
	}
	parseParams(r) // parse Params
	fetchInfo()    // fetch hitokoto info
	setResponse(w) // set content-type header
	setCallback(w) // check callback param
}
