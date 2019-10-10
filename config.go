package main

import (
	"io/ioutil"

	// "github.com/gomodule/redigo/redis"
	"github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v2"

	"fmt"
	"net/http"
	"os"
	"syscall"
)

// HITOKOTOAMOUNT is Number of databases
var HITOKOTOAMOUNT int64
var db *sqlx.DB
var err error
var conn *redis.Client

// FormatMap xxxxx
var FormatMap map[string]HTTPFormat
var config *UserConfig

type UserConfig struct {
	Postgres struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	}
	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	}
	ListenPort string `yaml:"listenport"`
}

// MakeReturnMap xxxxxxx
func MakeReturnMap() {
	FormatMap = make(map[string]HTTPFormat)
	FormatMap["js"] = HTTPFormat{Charset: "text/javascript; charset=", Text: "var hitokoto=\"%s——「%s」\";var dom=document.querySelector('.hitokoto');Array.isArray(dom)?dom[0].innerText=hitokoto:dom.innerText=hitokoto;"}
	FormatMap["json"] = HTTPFormat{Charset: "application/json; charset=", Text: "{\"hitokoto\": \"%s\", \"source\": \"%s\"}"}
	FormatMap["text"] = HTTPFormat{Charset: "text/plain; charset=", Text: "%s——「%s」"}
}

func initConfig(filename string) {
	config = new(UserConfig)
	yamlFile, _ := ioutil.ReadFile("./config.yaml")
	_ = yaml.Unmarshal(yamlFile, config)
}

func initRedis() {
	conn = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
}

func initHitokotoDB() {
	url := fmt.Sprintf("postgres://%s:%s@localhost/api?sslmode=disable", config.Postgres.User, config.Postgres.Password)
	db, err = sqlx.Connect("postgres", url)
	defer db.Close()
	checkErr(err)
}

func initLogFile() {
	file, _ := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
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

func IsLimited(r *http.Request) []interface{} {
	header := r.Header
	xforwared := header.Get("X-Forwarded-For")
	if xforwared == "" {
		xforwared = "NoForwaredIP"
	}
	ret, err := conn.Do("CL.THROTTLE", xforwared, "35", "36", "360").Result()
	checkErr(err)
	return ret.([]interface{})
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func setCheating(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprintf(w, "%s", "{'cheating': true}")
}
