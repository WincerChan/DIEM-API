package models

import (
	G "DIEM-API/models/googleanalytics"
	H "DIEM-API/models/hitokoto"

	bolt "go.etcd.io/bbolt"
)

func InitHitokoto(tx *bolt.Tx) error {
	H.HitokotoMapping = make(map[int]int)
	return H.ScanRecordLength(tx)
}

func InitGoogleAnalytics(viewID string) {
	InitGACredential()
	G.GAViewID = viewID
}

func InitMeiliSearch(host, apiKey string) {
	InitMeiliClient(host, apiKey)
}
