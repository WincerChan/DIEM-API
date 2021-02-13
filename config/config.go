package config

import (
	RPC "DIEM-API/rpcserver"
	T "DIEM-API/tools"
	"context"

	"github.com/go-redis/redis/v7"
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
	Pool         *RPC.Pool
	// GAViewID is same as before
	GAViewID string
	err      error
	// AnalyticsReportingService is same as before
	AnalyticsReportingService *gar.Service
)

// load config file(`config.yaml`) from disk.
func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	T.CheckFatalError(err, false)
	GAViewID = viper.GetString("google.analytics_id")
}

// init rpc server Connection-Pool
func initRPCServer() {
	Pool = RPC.NewPool(
		viper.GetInt("rpc-server.poolsize"),
		viper.GetString("rpc-server.addr"),
		RPC.DialTCP,
	)
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
	initRPCServer()
	initCredential()
}
