package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// HITOKOTOAMOUNT is Number of databases
var HITOKOTOAMOUNT int64
var db *sql.DB
var err error

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

func checkErr(err error) {
	switch {
	case err == sql.ErrNoRows:
		log.Printf(" None Query ")
	case err != nil:
		log.Fatal(err)
	}
}
