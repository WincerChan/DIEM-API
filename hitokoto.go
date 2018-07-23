package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
)

// Hit database composition
type Hit struct {
	Hitokoto string `json:"hitokoto"` // Hitokoto sentence
	Source   string `json:"source"`   // Hitokoto source
}

// Redirect301 old api
func Redirect301(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://api.itswincer.com/hitokoto/v2/", http.StatusMovedPermanently)
}

// Hitokoto handle function
func Hitokoto(w http.ResponseWriter, r *http.Request) {
	// query
	var hito string
	var source string
	var content string

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
		err1 := db.QueryRow("SELECT hitokoto, source FROM main WHERE LENGTH(hitokoto) < ? ORDER BY RAND() LiMIT 1;", length).Scan(&hito, &source)
		if err1 != nil {
			hito = ""
			source = ""
		}
	} else {
		nBig, err := rand.Int(rand.Reader, big.NewInt(HITOKOTOAMOUNT))
		n := nBig.Int64()
		err = db.QueryRow("SELECT hitokoto, source FROM main LIMIT ?, 1;", n).Scan(&hito, &source)
		checkErr(err)
	}

	if callback != "" {
		hs := &Hit{hito, source}
		fmtJSON, _ := json.Marshal(hs)
		w.Header().Set("Content-Type", "text/javascript; charset="+charset)
		fmt.Fprintf(w, "%s(%s);", callback, fmtJSON)
	} else {
		if encode == "js" {
			content = fmt.Sprintf("%s\\n——「%s」", hito, source)
			w.Header().Set("Content-Type", "text/javascript; charset="+charset)
			fmt.Fprintf(w, "var hitokoto=\"%s\";var dom=document.querySelector('.hitokoto');Array.isArray(dom)?dom[0].innerText=hitokoto:dom.innerText=hitokoto;", content)
		} else if encode == "json" {
			hs := &Hit{
				hito,
				source,
			}
			fmtJSON, _ := json.Marshal(hs)
			w.Header().Set("Content-Type", "application/json; charset="+charset)
			fmt.Fprintf(w, "%s", string(fmtJSON))
		} else if encode == "text" {
			w.Header().Set("Content-Type", "text/plain; charset="+charset)
			fmt.Fprintf(w, "%s", hito)
		} else {
			w.Header().Set("Content-Type", "text/plain; charset="+charset)
			content = fmt.Sprintf("%s——「%s」", hito, source)
			fmt.Fprintf(w, "%s", content)
		}
	}
}
