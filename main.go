package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
)

// AMOUNT is Number of databases
var AMOUNT int64
var db *sql.DB
var err error

// Config is pwd and user
type Config struct {
	User     string
	Password string
}

func initDatabase() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	url := fmt.Sprintf("%s:%s@/hitokoto?charset=utf8", config.User, config.Password)
	fmt.Println(config.User)
	db, err = sql.Open("mysql", url)
	checkErr(err)
	err1 := db.QueryRow("SELECT COUNT(id) FROM main;").Scan(&AMOUNT)
	checkErr(err1)
}

func main() {
	initDatabase()
	router := NewRouter()

	log.Println("listening in 520 port.")

	log.Fatal(http.ListenAndServe(":8020", router))
}

func checkErr(err error) {
	switch {
	case err == sql.ErrNoRows:
		log.Printf(" None Query ")
	case err != nil:
		log.Fatal(err)
	}
}
