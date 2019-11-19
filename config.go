package main

import (
	"io/ioutil"
	"strconv"

	// "github.com/gomodule/redigo/redis"
	"github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v2"

	"fmt"
	"net/http"
)

var db *sqlx.DB
var err error
var conn *redis.Client
var config *UserConfig

type UserConfig struct {
	Postgres struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	}
	Redis struct {
		Enabled  bool   `yaml:"enabled"`
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	}
	ListenPort string `yaml:"listenport"`
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
	url := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable",
		config.Postgres.User,
		config.Postgres.Password,
		config.Postgres.Database)
	db, err = sqlx.Connect("postgres", url)
	checkErr(err)
}

func getRemainingNumbers(r *http.Request) (limitInfo []string) {
	header := r.Header
	xforwared := header.Get("X-Forwarded-For")
	if xforwared == "" {
		xforwared = "NoForwaredIP"
	}
	ret, err := conn.Do("CL.THROTTLE", xforwared, "35", "36", "360").Result()
	checkErr(err)

	for _, v := range ret.([]interface{}) {
		v := v.(int64)                   // type assertion
		vStr := strconv.FormatInt(v, 10) // int64 to str
		limitInfo = append(limitInfo, vStr)
	}
	return
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
