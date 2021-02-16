package config

import (
	D "DIEM-API/models"
	RPC "DIEM-API/rpcserver"
	T "DIEM-API/tools"

	"github.com/spf13/viper"
	gar "google.golang.org/api/analyticsreporting/v4"
)

var (
	Pool     *RPC.Pool
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
	id := viper.GetString("google.analytics_id")
	D.InitGoogleAnalytics(id)
}

func initDatabase() {
	D.InitBoltConn(viper.GetString("bolt-path"))
	D.BoltDB.Read(D.InitHitokoto)
	D.InitMeiliSearch(viper.GetString("meilisearch.host"), viper.GetString("meilisearch.api-key"))
}

// init rpc server Connection-Pool
func initRPCServer() {
	Pool = RPC.NewPool(
		viper.GetInt("rpc-server.poolsize"),
		viper.GetString("rpc-server.addr"),
		RPC.DialTCP,
	)
}

// InitConfig init all config
func InitConfig() {
	loadConfig()
	initDatabase()
	initRPCServer()
}
