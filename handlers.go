package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// Hit database composition
type Hit struct {
	ID     string // Hitokoto ID
	HITO   string // Hitokoto sentence
	SOURCE string // Hitokoto source
}

// Hitokoto handle function
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
	} else if param == "word" {
		fmt.Fprintf(w, "%s", HITO)
	} else if param == "main" {
		fmt.Fprintf(w, "var hito = '%s\\n——「%s」'", HITO, SOURCE)
	} else {
		w.WriteHeader(404)
		fmt.Fprint(w, "error: Invalid API key")
	}
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
