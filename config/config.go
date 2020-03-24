package config

import (
	"fmt"
	"os"

	DNSTool "DIEM-API/tools/dnslookup"
	Logf "DIEM-API/tools/logfactory"

	"github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

var (
	RedisCli *redis.Client
	Enabled  bool
	PGConn   *sqlx.DB
	err      error
)

// init log configuration
func initLog() {
	Logf.Stderr = zerolog.ConsoleWriter{Out: os.Stderr}
	Logf.ErrLog = Logf.NewLogger("error")
	Logf.AccessLog = Logf.NewLogger("access")
	go Logf.ErrLog.Rotate()
	go Logf.AccessLog.Rotate()
}

// init redis connection
func initRedis() {
	address := DNSTool.ResolveAddr(viper.GetString("redis.address"))
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	_, err := RedisCli.Ping().Result()
	if err != nil {
		println("ERROR: Redis is not connected, disable rate-limiting.")
		return
	}
	_, err = RedisCli.Do("CL.THROTTLE", "", "35", "36", "360").Result()
	if err != nil {
		println("WARNING: No redis-cell module detected, disable rate-limiting.")
		return
	}
	if viper.GetBool("redis.enabled") {
		Enabled = true
	}
}

// load config file from disk.
func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		// TODO
		println("Load config failed.")
	}
}

// init PostgreSQL connection
func initPG() {
	host := DNSTool.ResolveOne(viper.GetString("postgres.host"))
	pgInfo := fmt.Sprintf(
		"host=%s user=%s port=%d dbname=%s sslmode=%s password=%s",
		host,
		viper.GetString("postgres.user"),
		viper.GetInt("postgres.port"),
		viper.GetString("postgres.database"),
		viper.GetString("postgres.sslmode"),
		viper.GetString("postgres.password"))
	PGConn, err = sqlx.Connect("postgres", pgInfo)
	if err != nil {
		panic(err)
	}
}

func init() {
	loadConfig()
	initLog()
	initPG()
	initRedis()
}
