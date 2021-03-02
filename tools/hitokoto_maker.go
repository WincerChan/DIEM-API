package tools

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
	"os"
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

var HitokotoMapping map[int]int

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

func insert(db *bolt.DB, id int, s string) {
	words := strings.Split(s, "\t")
	xxhash, err := strconv.Atoi(words[0])
	if err != nil {
		panic(err)
	}
	length, err := strconv.Atoi(words[2])
	if err != nil {
		panic(err)
	}
	r := &Record{
		Xxhash: int64(xxhash),
		Length: length,
		Origin: words[1],
	}
	r.Hitokoto.Source = words[3]
	r.Hitokoto.Hito = words[4]
	r.insert(db, uint32(id))
}

func OpenFile(path string) {
	db, _ := bolt.Open(path, 0666, nil)
	file, err := os.Open("./hito")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	defer db.Close()
	scanner := bufio.NewScanner(file)

	for i := 0; scanner.Scan(); i++ {
		insert(db, i, scanner.Text())
	}
}
