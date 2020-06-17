package config

import (
	T "DIEM-API/tools"
	DNSTool "DIEM-API/tools/dnslookup"
	"context"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	gar "google.golang.org/api/analyticsreporting/v4"
	"google.golang.org/api/option"

	// postgres library
	_ "github.com/lib/pq"
)

var (
	// RedisCli is in package scope variable
	RedisCli *redis.Client
	// EnabledRedis is same as before
	EnabledRedis bool
	// PGConn is same as before
	PGConn *sqlx.DB
	// GAViewID is same as before
	GAViewID string
	err      error
	// AnalyticsReportingService is same as before
	AnalyticsReportingService *gar.Service
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

// load config file(`config.yaml`) from disk.
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

// init Google Analytics credential
func initCredential() {
	ctx := context.Background()
	json := T.LoadJSON("./credential.json")
	AnalyticsReportingService, err = gar.NewService(ctx, option.WithCredentialsJSON(json))
	T.CheckFatalError(err, false)
}

// InitConfig init all config
func InitConfig() {
	loadConfig()
	initPG()
	initRedis()
	initCredential()
}
