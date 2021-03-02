package config

import (
	D "DIEM-API/models"
	RPC "DIEM-API/rpcserver"
	T "DIEM-API/tools"
	L "DIEM-API/tools/logfactory"
	"os"

	"github.com/spf13/viper"
	gar "google.golang.org/api/analyticsreporting/v4"
)

var (
	RalPool    *RPC.Pool
	SearchPool *RPC.Pool
	GAViewID   string
	err        error
	// AnalyticsReportingService is same as before
	AnalyticsReportingService *gar.Service
)

// load config file(`config.yaml`) from disk.
func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(os.Getenv("CONFIG"))
	err = viper.ReadInConfig()
	T.CheckFatalError(err, false)
	id := viper.GetString("google.analytics_id")
	D.InitGoogleAnalytics(id)
	L.InitLog(T.ConfigAbsPath("_logs"))
}

func initDatabase() {
	D.InitBoltConn(T.ConfigAbsPath(viper.GetString("bolt-path")))
	D.BoltDB.Read(D.InitHitokoto)
}

// init rpc server Connection-Pool
func initRPCServer() {
	RalPool = RPC.NewPool(
		viper.GetInt("rpc-server.poolsize"),
		T.ConfigAbsPath(viper.GetString("rpc-server.addr")),
		RPC.DialTCP,
	)
	SearchPool = RPC.NewPool(
		viper.GetInt("search.poolsize"),
		viper.GetString("search.addr"),
		RPC.DialTCP,
	)
}

// InitConfig init all config
func InitConfig() {
	loadConfig()
	initDatabase()
	initRPCServer()
}
