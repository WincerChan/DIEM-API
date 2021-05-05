package rpcserver

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
)

func BenchmarkToString(b *testing.B) {
	var bs []byte
	for i := 0; i < b.N; i++ {
		s := strconv.Itoa(37)
		bs = append(bs, []byte(s)...)
		f := fmt.Sprintf("%f", 0.3497)
		bs = append(bs, []byte(f)...)
		i := strconv.Itoa(37)
		bs = append(bs, []byte(i)...)
		bs = []byte{}
	}
}

func BenchmarkTLVEncode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		bf := new(bytes.Buffer)
		encodeString(bf, "choke")
		encodeInteger(bf, 37)
		encodeFloat(bf, 0.3498)
		bf.Bytes()
	}
}
