package config

import (
	"bytes"
	"encoding/gob"
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

func (r *Record) LoadFromBytes(value []byte) {
	buf := new(bytes.Buffer)
	buf.Write(value)
	gob.NewDecoder(buf).Decode(r)
}
