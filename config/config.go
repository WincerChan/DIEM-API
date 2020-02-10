package config

import (
	"fmt"
	"os"

	Logf "DIEM-API/tools/logfactory"

	"github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

var (
	RedisCli *redis.Client
	PGConn   *sqlx.DB
	err      error
)

func initLog() {
	Logf.Stderr = zerolog.ConsoleWriter{Out: os.Stderr}
	Logf.ErrLog = Logf.NewLogger("error")
	Logf.AccessLog = Logf.NewLogger("access")
	go Logf.ErrLog.Rotate()
	go Logf.AccessLog.Rotate()
}

func initRedis() {
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.address"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	_, err := RedisCli.Ping().Result()
	if err != nil {
		println("ERROR: Sorry, Redis is not connected.")
		os.Exit(1)
	}
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		// TODO
		println("Load config failed.")
	}
}

func initPG() {
	pgInfo := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		viper.GetString("postgres.user"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.host"),
		viper.GetInt("postgres.port"),
		viper.GetString("postgres.database"),
		viper.GetString("postgres.sslmode"))
	PGConn, err = sqlx.Connect("postgres", pgInfo)
	if err != nil {
		panic(err)
	}
}

func init() {
	initLog()
	initConfig()
	initPG()
	initRedis()
}
