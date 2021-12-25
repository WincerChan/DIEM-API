package config

import (
	D "DIEM-API/models"
	RPC "DIEM-API/rpcserver"
	T "DIEM-API/tools"
	L "DIEM-API/tools/logfactory"
	C "DIEM-API/tools/tomlparser"
	V "DIEM-API/views"

	R "DIEM-API/middleware/limiting"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	gar "google.golang.org/api/analyticsreporting/v4"
)

var (
	configPath string
	GAViewID   string
	// AnalyticsReportingService is same as before
	AnalyticsReportingService *gar.Service
	RegisterService           []string
)

func initLogService() {
	L.InitLog()
}

// load config file(`config.yaml`) from disk.
func loadConfig() {
	id := C.GetString("credential.analytics-id")
	credentialPath := C.ConfigAbsPath("credential.filename")
	D.InitGoogleAnalytics(id, credentialPath)
}

func initDatabase() {
	path := C.ConfigAbsPath("hitokoto.dbpath")
	D.InitBoltConn(path)
	D.BoltDB.Read(D.InitHitokoto)
}

// init rpc server Connection-Pool
func initRPCServer() {
	var addr string
	if C.GetString("rate-limit.network") == "uds" {
		addr = C.ConfigAbsPath("rate-limit.addr")
	} else {
		addr = C.GetString("rate-limit.addr")
	}
	R.InitRateLimit("uds", addr, C.GetInt("rate-limit.poolsize"))
}

func initSearchAPI() {
	var addr string
	if C.GetString("search.network") == "uds" {
		addr = C.ConfigAbsPath("search.addr")
	} else {
		addr = C.GetString("search.addr")
	}
	V.SearchPool = RPC.NewPool(
		C.GetInt("search.poolsize"),
		addr,
		RPC.DialTCP,
	)
}

// InitConfig init all config
func InitConfig(conf string) {
	configPath = conf
}
func InitService(r *gin.Engine, service string) {
	C.LoadTOML(configPath)
	initRPCServer()
	initLogService()
	if strings.HasPrefix("hitokoto", service) {
		initDatabase()
	}
	if strings.HasPrefix("analytics", service) {
		loadConfig()
	}
	if strings.HasPrefix("search", service) {
		initSearchAPI()
	}
	V.Register(r, service)
}

func MigrateBolt() {
	C.LoadTOML(configPath)
	path := C.ConfigAbsPath("hitokoto.dbpath")
	os.Remove(path)
	println("Trying to migrate database")
	T.MigrateHitokoto(C.ConfigAbsPath("hitokoto.source"), path)
	println("succeed.")
}
