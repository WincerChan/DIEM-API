package config

import (
	T "DIEM-API/tools"
	DNSTool "DIEM-API/tools/dnslookup"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

var (
	RedisCli     *redis.Client
	EnabledRedis bool
	PGConn       *sqlx.DB
	GAViewID     string
	err          error
)

// init redis connection
func initRedis() {
	if !viper.GetBool("redis.enabled") {
		return
	}
	address := DNSTool.ResolveAddr(viper.GetString("redis.address"))
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	_, err = RedisCli.Ping().Result()
	T.CheckException(err,
		"WARNING: Couldn't connect to Redis, please set redis.enabled to false.")
	_, err = RedisCli.Do("CL.THROTTLE", "", "35", "36", "360").Result()
	T.CheckException(err,
		"WARNING: redis-cell module didn't detect, please set redis.enabled to false.")
	EnabledRedis = true
}

// load config file from disk.
func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	T.CheckFatalError(err, false)
	GAViewID = viper.GetString("google.analytics_id")
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
	T.CheckFatalError(err, false)
}

func init() {
	loadConfig()
	initPG()
	initRedis()
}
