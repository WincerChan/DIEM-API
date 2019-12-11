package config

import (
	"fmt"
	"io"
	"os"

	T "DIEM-API/tools"

	"github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

// LogOutput
var LogOutput struct {
	stderr    io.Writer
	errLog    *T.Logger
	accessLog *T.Logger
}

var (
	RedisCli *redis.Client
	PGConn   *sqlx.DB
	err      error
)

func initLog() {
	LogOutput.stderr = zerolog.ConsoleWriter{Out: os.Stderr}
	LogOutput.errLog = T.NewLogger("error")
	LogOutput.accessLog = T.NewLogger("access")
	go LogOutput.errLog.Rotate()
	go LogOutput.accessLog.Rotate()
}

func initRedis() {
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.address"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	// println("Addr", viper.GetString("redis.address"), viper.GetString("redis.password"), viper.GetInt("redis.db"))
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
	println(pgInfo)
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

// 应该放在 tools 方便取出来用
func GetWriter(w string) io.Writer {
	switch w {
	case "std":
		return LogOutput.stderr
	case "error":
		return LogOutput.errLog.Writer
	case "access":
		return LogOutput.accessLog.Writer
	}
	return nil
}
