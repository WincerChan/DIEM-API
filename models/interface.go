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

func InitGoogleAnalytics(viewID, filepath string) {
	InitGACredential(filepath)
	G.GAViewID = viewID
}
