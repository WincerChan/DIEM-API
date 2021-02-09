package hitokoto

import (
	C "DIEM-API/config"
	T "DIEM-API/tools"
	"math/rand"

	bolt "go.etcd.io/bbolt"
)

// MAXRECORD number of hitokoto
const MAXRECORD = 15348

func getOffset(length int) (offset int) {
	if val, ok := C.HitokotoMapping[length]; ok {
		offset = rand.Intn(val)
	} else {
		offset = rand.Intn(MAXRECORD)
	}
	return
}

// fetch hitokoto from database
func FetchHitokoto(length int) *C.HitoInfo {
	key, record := T.Int32ToBytes(getOffset(length)), new(C.Record)
	C.BoltConn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("hitokoto"))
		record.LoadFromBytes(b.Get(key))
		return nil
	})
	return &record.Hitokoto
}
