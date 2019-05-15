package main

// This file has been discarded
// This file has been discarded
// This file has been discarded
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Thinking struct {
	ID       string    `db:"id"`
	Created  time.Time `db:"created"`
	Text     string    `db:"text"`
	Position string    `db:"position"`
	LikeNum  int       `db:"like_num"`
	Liked    bool      `db:"liked"`
}

type Tid struct {
	ID string `json:"tid"`
}

func setReqHeader(r *http.Request) {
	headers := r.Header
	cfip := headers.Get("Cf-Connecting-Ip")
	xff := headers.Get("X-Forwarded-For")
	cfco := headers.Get("Cf-Ipcountry")
	ref := headers.Get("Referer")
	ua := headers.Get("User-Agent")
	co := headers.Get("Cookie")
	path := r.Method + r.URL.Path
	_, err := reqDB.Exec(`INSERT req SET 
	cfconnectingip=?, xforwardedfor=?, 
	cfipcountry=?, referer=?, 
	cookie=?, useragent=?, urlpath=?`, cfip, xff, cfco, ref, co, ua, path)
	checkErr(err)
}

func checkCfuidValid(cfuid string) bool {
	if len(cfuid) != 43 {
		return false
	}
	return true
}

func queryThinking(from, size int, cfuid string) []byte {
	Query := `
	SELECT t1.*, IFNULL(t2.liked, FALSE) as liked, IFNULL(t3.like_num, 0) as like_num
	FROM  
	(SELECT * 
	FROM thinking ORDER BY created DESC
	limit ?, ?) AS t1
	LEFT JOIN  
	(SELECT thinking, True AS liked 
	FROM thinking_user 
	where user =?) AS t2 
	ON t1.id=t2.thinking 
	LEFT JOIN 
	(SELECT thinking, COUNT(thinking) AS like_num 
	FROM thinking_user 
	GROUP BY thinking) AS t3 
	ON t1.id = t3.thinking
	`
	ts := []Thinking{}
	err := tkDB.Select(&ts, Query, from, size, cfuid)
	checkErr(err)
	fmtJSON, _ := json.Marshal(ts)
	return fmtJSON
}

func insertThinking(tid, cfuid string) {
	log.Printf("tid=%s, cfuid=%s\n", tid, cfuid)
	tkx, err := tkDB.Begin()
	_, err = tkx.Exec("INSERT INTO thinking_user (thinking, user) VALUES (?, ?);", tid, cfuid)
	if err != nil {
		_, err := tkx.Exec("INSERT INTO user (id) values (?)", cfuid)
		checkErr(err)
		_, err = tkx.Exec("INSERT INTO thinking_user (thinking, user) VALUES (?, ?);", tid, cfuid)
	}
	err = tkx.Commit()
	if err != nil {
		tkx.Rollback()
	}
}

func getThinkMethod(w http.ResponseWriter, r *http.Request, cfuid string) []byte {
	r.ParseForm()
	size := 10
	from, err := strconv.Atoi(r.Form.Get("from"))
	if err != nil {
		setCheating(w)
		return []byte("")
	}
	result := queryThinking(from, size, cfuid)
	return result

}

func postThinkMethod(w http.ResponseWriter, r *http.Request, cfuid string) []byte {
	tid := Tid{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	_ = json.Unmarshal(body, &tid)
	insertThinking(tid.ID, cfuid)
	return body
}

func HandleThinkReq(w http.ResponseWriter, r *http.Request) {
	setReqHeader(r)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	cfuid, err := r.Cookie("__cfduid")
	if err != nil || !checkCfuidValid(cfuid.Value) {
		setCheating(w)
		return
	}
	cfuidStr := cfuid.Value
	var queryResult []byte

	switch r.Method {
	case http.MethodGet:
		queryResult = getThinkMethod(w, r, cfuidStr)
	case http.MethodPost:
		queryResult = postThinkMethod(w, r, cfuidStr)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	fmt.Fprintf(w, "%s", string(queryResult))
}
