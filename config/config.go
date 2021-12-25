package config

import (
	D "DIEM-API/models"
	L "DIEM-API/tools/logfactory"
	C "DIEM-API/tools/tomlparser"
	V "DIEM-API/views"

	R "DIEM-API/middleware/limiting"
	"strings"

	"github.com/gin-gonic/gin"
)

func initLogService() {
	L.InitLog()
}

// load config file(`config.yaml`) from disk.
func loadCredential() {
	id := C.GetString("credential.analytics-id")
	credentialPath := C.ConfigAbsPath("credential.filename")
	D.InitGoogleAnalytics(id, credentialPath)
}

func loadAddrFromConfig(component string) (net, addr string) {
	networkType := component + ".network"
	addrPath := component + ".addr"
	net = C.GetString(networkType)
	if net == "uds" {
		addr = C.ConfigAbsPath(addrPath)
	} else {
		addr = C.GetString(addrPath)
	}
	return
}

func initDatabase() {
	path := C.ConfigAbsPath("hitokoto.dbpath")
	D.InitBoltConn(path)
}

// init rpc server Connection-Pool
func initRPCServer() {
	net, addr := loadAddrFromConfig("rate-limit")
	R.InitRalPool(net, addr, C.GetInt("rate-limit.poolsize"))
}

func initSearchAPI() {
	net, addr := loadAddrFromConfig("search")
	V.InitSearchPool(net, addr, C.GetInt("search.poolsize"))
}

// InitConfig init all config
func InitConfig(conf string) {
	C.LoadTOML(conf)
	initCommonService()
}

func initCommonService() {
	initRPCServer()
	initLogService()
}

func InitService(r *gin.Engine, service string) {
	if strings.HasPrefix("hitokoto", service) {
		initDatabase()
	}
	if strings.HasPrefix("analytics", service) {
		loadCredential()
	}
	if strings.HasPrefix("search", service) {
		initSearchAPI()
	}
	V.Register(r, service)
}
