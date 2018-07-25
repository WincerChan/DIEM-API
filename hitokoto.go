package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
)

// Hit database composition
type Hit struct {
	Hitokoto string `json:"hitokoto"` // Hitokoto sentence
	Source   string `json:"source"`   // Hitokoto source
}

// query
var hito string
var source string
var content string

// Redirect301 old api
func Redirect301(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://api.itswincer.com/hitokoto/v2/", http.StatusMovedPermanently)
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

// HasLength if url param has length
func HasLength(length string) {
	lengthInt, err := strconv.Atoi(length)
	if err != nil || lengthInt < 5 {
		handleError("invalidLengthError")
		return
	}
	err1 := db.QueryRow("SELECT hitokoto, source FROM main WHERE LENGTH(hitokoto) < ? ORDER BY RAND() LiMIT 1;", length).Scan(&hito, &source)
	checkErr(err1)
}

// GenRandomInt to gen a random int
func GenRandomInt() {
	nBig, err := rand.Int(rand.Reader, big.NewInt(HITOKOTOAMOUNT))
	n := nBig.Int64()
	err = db.QueryRow("SELECT hitokoto, source FROM main LIMIT ?, 1;", n).Scan(&hito, &source)
	checkErr(err)
}

// Hitokoto handle function
func Hitokoto(w http.ResponseWriter, r *http.Request) {
	if DisallowMethod(w, "GET", r.Method) {
		// if method not allow, just return
		return
	}

	// get params
	r.ParseForm()
	encode := r.Form.Get("encode")
	length := r.Form.Get("length")
	// if url param have callback then will ignore encode
	callback := r.Form.Get("callback")
	charset := r.Form.Get("charset")
	if charset == "gbk" {
	} else {
		charset = "utf-8"
	}
	log.Println(r.URL.Path)

	if length != "" {
		HasLength(length)
	} else {
		GenRandomInt()
	}
	// hasCallback is return data
	hasCallback := func() {
		if callback != "" {
			hs := &Hit{hito, source}
			fmtJSON, _ := json.Marshal(hs)
			w.Header().Set("Content-Type", "text/javascript; charset="+charset)
			fmt.Fprintf(w, "%s(%s);", callback, fmtJSON)
			return
		}
		switch encode {
		case "js":
			content = fmt.Sprintf("%s\\n——「%s」", hito, source)
			w.Header().Set("Content-Type", "text/javascript; charset="+charset)
			fmt.Fprintf(w, "var hitokoto=\"%s\";var dom=document.querySelector('.hitokoto');Array.isArray(dom)?dom[0].innerText=hitokoto:dom.innerText=hitokoto;", content)
		case "json":
			hs := &Hit{
				hito,
				source,
			}
			fmtJSON, _ := json.Marshal(hs)
			w.Header().Set("Content-Type", "application/json; charset="+charset)
			fmt.Fprintf(w, "%s", string(fmtJSON))
		case "text":
			w.Header().Set("Content-Type", "text/plain; charset="+charset)
			fmt.Fprintf(w, "%s", hito)
		default:
			w.Header().Set("Content-Type", "text/plain; charset="+charset)
			content = fmt.Sprintf("%s——「%s」", hito, source)
			fmt.Fprintf(w, "%s", content)
		}
	}
	hasCallback()
}
