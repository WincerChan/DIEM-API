package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/jmoiron/sqlx"
)

// HITOKOTOAMOUNT is Number of databases
var HITOKOTOAMOUNT int64
var db, tkDB, reqDB *sqlx.DB
var err error
var file *os.File
var config MysqlCONF

// FormatMap xxxxx
var FormatMap map[string]HTTPFormat

// MysqlCONF is pwd and user
type MysqlCONF struct {
	User     string
	Password string
	Port     string
}

// MakeReturnMap xxxxxxx
func MakeReturnMap() {
	FormatMap = make(map[string]HTTPFormat)
	FormatMap["js"] = HTTPFormat{Charset: "text/javascript; charset=", Text: "var hitokoto=\"%s——「%s」\";var dom=document.querySelector('.hitokoto');Array.isArray(dom)?dom[0].innerText=hitokoto:dom.innerText=hitokoto;"}
	FormatMap["json"] = HTTPFormat{Charset: "application/json; charset=", Text: "{\"hitokoto\": \"%s\", \"source\": \"%s\"}"}
	FormatMap["text"] = HTTPFormat{Charset: "text/plain; charset=", Text: "%s——「%s」"}
}
func initHitokotoDB(filename string) {
	file, _ := os.Open(filename)
	decoder := json.NewDecoder(file)
	config = MysqlCONF{}
	err := decoder.Decode(&config)
	url := fmt.Sprintf("%s:%s@/hitokoto?charset=utf8", config.User, config.Password)
	tkURL := fmt.Sprintf("%s:%s@/THINKS?charset=utf8&parseTime=true", config.User, config.Password)
	reqURL := fmt.Sprintf("%s:%s@/apidata?charset=utf8", config.User, config.Password)
	// db, err = sql.Open("mysql", url)
	db, err = sqlx.Connect("mysql", url)
	tkDB, err = sqlx.Connect("mysql", tkURL)
	reqDB, err = sqlx.Connect("mysql", reqURL)
	// fmt.Println(reflect.TypeOf(b))
	checkErr(err)
	err1 := db.QueryRow("SELECT COUNT(id) FROM main;").Scan(&HITOKOTOAMOUNT)
	checkErr(err1)
}

func initLogFile() {
	file, err = os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	defer file.Close()
	syscall.Dup2(int(file.Fd()), 1)
	syscall.Dup2(int(file.Fd()), 2)
}

//DisallowMethod is allow current method
func DisallowMethod(w http.ResponseWriter, allow string, method string) bool {
	if allow != method && allow != "HEAD" {
		w.WriteHeader(405)
		fmt.Fprintln(w, "<h1>405 Not Allowed</h1>")
		return true
	}
	return false
}

func checkErr(err error) {
	switch {
	case err == sql.ErrNoRows:
		// handleError("queryError")
		hito = "哦~"
		source = "袴田日向"
		log.Println("None Query")
	case err != nil:
		// handleError("connectError")
		hito = "哦~"
		source = "袴田日向"
		log.Println(err)
	}
}

func setCheating(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprintf(w, "%s", "{'cheating': true}")
}
