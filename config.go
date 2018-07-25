package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"
)

// HITOKOTOAMOUNT is Number of databases
var HITOKOTOAMOUNT int64
var db *sql.DB
var err error
var file *os.File

// MysqlCONF is pwd and user
type MysqlCONF struct {
	User     string
	Password string
}

func initHitokotoDB() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config := MysqlCONF{}
	err := decoder.Decode(&config)
	url := fmt.Sprintf("%s:%s@/hitokoto?charset=utf8", config.User, config.Password)
	db, err = sql.Open("mysql", url)
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
	if allow != method {
		w.WriteHeader(405)
		fmt.Fprintln(w, "<h1>405 Not Allowed</h1>")
		return true
	}
	return false
}

func checkErr(err error) {
	switch {
	case err == sql.ErrNoRows:
		handleError("queryError")
		log.Println("None Query")
	case err != nil:
		handleError("connectError")
		log.Println(err)
	}
}
