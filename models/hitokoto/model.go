package hitokoto

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"

	bolt "go.etcd.io/bbolt"
)

var (
	counts          int
	HitokotoMapping map[int]int
	HitoBucket      = []byte("hitokoto")
)

type Params struct {
	Length   int    `form:"length"`
	Callback string `form:"callback"`
	Encode   string `form:"encode"`
}

type HitoInfo struct {
	Source string `json:"source"`
	Hito   string `json:"hitokoto"`
}

type Record struct {
	Xxhash   int64
	Length   int
	Origin   string
	Hitokoto HitoInfo
}

func LoadRecordFromBytes(value []byte) Record {
	buf, r := new(bytes.Buffer), new(Record)
	buf.Write(value)
	gob.NewDecoder(buf).Decode(r)
	return *r
}

func IndexOf(length int) int {
	if val, ok := HitokotoMapping[length]; ok {
		return val
	}
	return counts
}

func ScanRecordLength(tx *bolt.Tx) error {
	b := tx.Bucket(HitoBucket)
	preLength := 0
	b.ForEach(func(k, v []byte) error {
		counts++
		id := binary.BigEndian.Uint32(k)
		record := LoadRecordFromBytes(v)
		if record.Length != preLength {
			HitokotoMapping[record.Length] = int(id)
		}
		preLength = record.Length
		return nil
	})
	return nil
}
