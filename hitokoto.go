package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

// HTTPFormat xxxx
type HTTPFormat struct {
	Charset string `json:"charset"`
	Text    string `json:"text"`
}

type Content struct {
	Source string `json:"source"`
	Hito   string `json:"hitokoto"`
}

func (c Content) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Content) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}

// query
var cnt *Content
var hito string
var source string
var content string

// Hitokoto handle function
func Hitokoto(w http.ResponseWriter, r *http.Request) {
	ret := IsLimited(r)
	w.Header().Set("X-RateLimit-Limit", strconv.FormatInt(ret[1].(int64), 10))
	w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(ret[2].(int64), 10))
	w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(ret[4].(int64), 10))
	if ret[0].(int64) == 1 {
		content = "{\"result\": \"Your IP requests is frequently.\"}"
	} else {
		cnt := new(Content)
		log.Println(r.URL.Path)
		// setReqHeader(r)
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
		if length == "" || len(length) > 3 {
			db.QueryRow("SELECT RANDOMFETCH($1);", -1).Scan(&cnt)
		} else {
			lengthInt, err := strconv.Atoi(length)
			if err != nil {
				checkErr(err)
			} else {
				db.QueryRow("SELECT RANDOMFETCH($1);", lengthInt).Scan(&cnt)
			}
		}
		cntJSON, _ := cnt.Value()
		// hasCallback is return data
		w.Header().Set("Content-Type", FormatMap["text"].Charset+charset)
		// The value that needs to be returned
		content = fmt.Sprintf(FormatMap["text"].Text, cnt.Hito, cnt.Source)
		// set content to encode format
		if text, ok := FormatMap[encode]; ok {
			w.Header().Set("Content-Type", text.Charset+charset)
			content = fmt.Sprintf(text.Text, cnt.Hito, cnt.Source)
		}
		// if url params have callback then will ignore encode
		if callback != "" {
			w.Header().Set("Content-Type", "text/javascript; charset="+charset)
			content = fmt.Sprintf("%s(%s)", callback, string(cntJSON.([]byte)))
		}

	}
	// output content
	fmt.Fprint(w, content)

}
