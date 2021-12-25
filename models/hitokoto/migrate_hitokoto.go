package hitokoto

import (
	T "DIEM-API/tools"
	L "DIEM-API/tools/logfactory"
	C "DIEM-API/tools/tomlparser"
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	bolt "go.etcd.io/bbolt"
)

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

type SortBy []Record

func (a SortBy) Len() int           { return len(a) }
func (a SortBy) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortBy) Less(i, j int) bool { return a[i].Length < a[j].Length }

func (r *Record) insert(db *bolt.DB, id uint32) {
	key := make([]byte, 4)
	binary.BigEndian.PutUint32(key, id)
	value := new(bytes.Buffer)
	err := gob.NewEncoder(value).Encode(r)
	if err != nil {
		log.Println(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("hitokoto"))
		if err != nil {
			return err
		}
		return b.Put(key, value.Bytes())
	})
}

func formatAsRecord(data []string) []Record {
	records := make([]Record, len(data), len(data))
	for i, text := range data {
		words := strings.Split(text, "\t")
		xxhash, err := strconv.Atoi(words[0])
		T.CheckFatalError(err, false)
		length, err := strconv.Atoi(words[2])
		T.CheckFatalError(err, false)
		r := &Record{
			Xxhash: int64(xxhash),
			Length: length,
			Origin: words[1],
		}
		r.Hitokoto.Source = words[3]
		r.Hitokoto.Hito = words[4]
		records[i] = *r
	}
	return records
}

func bulkInsert(db *bolt.DB, data []string) {
	records := formatAsRecord(data)
	sort.Sort(SortBy(records))
	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("hitokoto"))
		T.CheckFatalError(err, false)
		for i, r := range records {
			key := make([]byte, 4)
			binary.BigEndian.PutUint32(key, uint32(i))
			value := new(bytes.Buffer)
			err := gob.NewEncoder(value).Encode(r)
			T.CheckFatalError(err, false)
			b.Put(key, value.Bytes())
		}
		return nil
	})
}

func migrateHitokoto(source, path string) {
	db, _ := bolt.Open(path, 0666, nil)
	file, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	defer db.Close()
	hitokotos := make([]string, 0)
	scanner := bufio.NewScanner(file)

	for i := 0; scanner.Scan(); i++ {
		hitokotos = append(hitokotos, scanner.Text())
	}
	bulkInsert(db, hitokotos)
}

func MigrateBolt() {
	path := C.ConfigAbsPath("hitokoto.dbpath")
	source := C.ConfigAbsPath("hitokoto.source")
	os.Remove(path)
	L.Error.Debug().Msg("Trying to migrate database")
	migrateHitokoto(source, path)
	L.Error.Debug().Msg("Succeed.")
}
