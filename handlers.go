package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
)

type Hit struct {
	ID     string
	HITO   string
	SOURCE string
}

func Hitokoto(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "HITODB.db")
	checkErr(err)

	// query
	var ID string
	var HITO string
	var SOURCE string

	err1 := db.QueryRow("SELECT * FROM hitokoto ORDER BY RANDOM() LIMIT 1").Scan(&ID, &HITO, &SOURCE)
	checkErr(err1)

	// get params
	param := r.URL.Query().Get("encode")
	if param == "js" {
		fmt.Fprintf(w, "function hitokoto(){document.write('%s&#10;——「%s」');}", HITO, SOURCE)
	} else if param == "json" {
		hh := &Hit{ID, HITO, SOURCE}
		js, _ := json.Marshal(hh)
		fmt.Fprintf(w, "%s", js)
	} else if param == "main" {
		fmt.Fprintf(w, "%s", HITO)
	}
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func openLogFile(logfile string) {
	if logfile != "" {
		lf, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)

		if err != nil {
			log.Fatal("OpenLogfile: os.Openfile:", err)
		}

		log.SetOutput(lf)
	}
}
