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

var hitoinfo *Info
var done bool

type Info struct {
	Source string `json:"source"`
	Hito   string `json:"hitokoto"`
}

func (c Info) Value() []byte {
	result, _ := json.Marshal(c)
	return result
}

func (c *Info) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}

func setLimitHeader(w http.ResponseWriter, r *http.Request) (limited bool) {
	// If user do not enable redis-cell limit, just do not check
	if !config.Redis.Enabled {
		return
	}

	ret := getRemainingNumbers(r)
	w.Header().Set("X-RateLimit-Limit", ret[1])
	w.Header().Set("X-RateLimit-Remaining", ret[2])
	w.Header().Set("X-RateLimit-Reset", ret[4])

	if ret[0] == "0" {
		return
	}

	w.Write([]byte("Sorry, Your IP requests is too frequently."))
	return true
}

func parseParams(r *http.Request, pipe chan string) {
	r.ParseForm()
	pipe <- r.Form.Get("length")
	pipe <- r.Form.Get("callback")
	pipe <- r.Form.Get("charset")
	pipe <- r.Form.Get("encode")
}

func fetchInfo(pipe chan string) {
	lenStr := <-pipe
	if lenStr == "" || len(lenStr) > 3 {
		db.QueryRow("SELECT RANDOMFETCH($1);", -1).Scan(&hitoinfo)
	} else {
		length, err := strconv.Atoi(lenStr)
		checkErr(err)
		db.QueryRow("SELECT RANDOMFETCH($1);", length).Scan(&hitoinfo)
	}
}

func makeSTDResponse(w http.ResponseWriter, pipe chan string) {
	if done {
		return
	}

	charset := <-pipe
	if "gbk" != charset {
		charset = "utf-8"
	}

	var buffer bytes.Buffer
	var contenttype string

	switch e := <-pipe; e {
	case "js":
		contenttype = "text/javascript; charset=" + charset
		buffer.WriteString("var hitokoto=\"")
		buffer.WriteString(hitoinfo.Hito)
		buffer.WriteString("——「")
		buffer.WriteString(hitoinfo.Source)
		buffer.WriteString("」\";var dom=document.querySelector('.hitokoto');")
		buffer.WriteString("Array.isArray(dom)?dom[0].innerText=hitokoto:dom.innerText=hitokoto;")
	case "json":
		contenttype = "application/json; charset=" + charset
		buffer.WriteString("{\"hitokoto\": \"")
		buffer.WriteString(hitoinfo.Hito)
		buffer.WriteString("\", \"source\": \"")
		buffer.WriteString(hitoinfo.Source)
		buffer.WriteString("\"}")
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

func makeCallback(w http.ResponseWriter, pipe chan string) {
	callback := <-pipe
	if "" == callback {
		return
	}

	w.Header().Set("Content-Type", "text/javascript")

	var buffer bytes.Buffer
	buffer.WriteString(callback)
	buffer.WriteString("(")
	buffer.WriteString(string(hitoinfo.Value()))
	buffer.WriteString(")")
	w.Write(buffer.Bytes())
	done = true
}

// Hitokoto handle function
func Hitokoto(w http.ResponseWriter, r *http.Request) {
	pipe := make(chan string, 4)
	done = false
	log.Println(r.URL.Path)
	isLimited := setLimitHeader(w, r)
	if isLimited {
		return
	}
	parseParams(r, pipe)     // parse Params
	fetchInfo(pipe)          // fetch hitokoto info
	makeCallback(w, pipe)    // check callback param
	makeSTDResponse(w, pipe) // make standard response
	close(pipe)
}
