package config

import (
	RPC "DIEM-API/rpcserver"
	T "DIEM-API/tools"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/gob"

	"github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
	gar "google.golang.org/api/analyticsreporting/v4"
	"google.golang.org/api/option"

	bolt "go.etcd.io/bbolt"
)

var (
	// RedisCli is in package scope variable
	RedisCli *redis.Client
	// EnabledRedis is same as before
	EnabledRedis bool
	// BoltConn is same as before
	BoltConn *bolt.DB
	Pool     *RPC.Pool
	// GAViewID is same as before
	GAViewID string
	err      error
	// AnalyticsReportingService is same as before
	AnalyticsReportingService *gar.Service
	HitokotoMapping           map[int]int
)

// load config file(`config.yaml`) from disk.
func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	T.CheckFatalError(err, false)
	GAViewID = viper.GetString("google.analytics_id")
}

func initBolt() {
	HitokotoMapping = make(map[int]int)
	BoltConn, err = bolt.Open("/tmp/bbolt", 0666, &bolt.Options{ReadOnly: true})
	T.CheckFatalError(err, false)
	BoltConn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("hitokoto"))
		preLength := 0
		b.ForEach(func(k, v []byte) error {
			id := binary.BigEndian.Uint32(k)
			buf, r := new(bytes.Buffer), new(Record)
			buf.Write(v)
			gob.NewDecoder(buf).Decode(&r)
			if r.Length != preLength {
				HitokotoMapping[r.Length] = int(id)
			}
			preLength = r.Length
			return nil
		})
		return nil
	})
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
	initBolt()
	initRPCServer()
	initCredential()
}
