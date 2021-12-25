package models

import (
	T "DIEM-API/tools"

	bolt "go.etcd.io/bbolt"
)

type BoltConn struct{ bolt.DB }

var BoltDB *BoltConn

func (bc *BoltConn) Read(fn func(tx *bolt.Tx) error) {
	bc.View(fn)
}

func (bc *BoltConn) Write(fn func(tx *bolt.Tx) error) {
	bc.Update(fn)
}

func InitBoltConn(path string) {
	db, err := bolt.Open(path, 0666, &bolt.Options{ReadOnly: true})
	T.CheckFatalError(err, false)
	BoltDB = &BoltConn{*db}
	BoltDB.Read(InitHitokoto)
}
