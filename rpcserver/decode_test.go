package rpcserver

import (
	"strconv"
	"strings"
	"testing"
)

func BenchmarkStringTo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		a := "choke$37$0.3498"
		strs := strings.Split(a, "$")
		_ = strs[0]
		strconv.Atoi(strs[1])
		strconv.ParseFloat(strs[2], 64)
	}
}

func BenchmarkTLVDecode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		bs := []byte{0, 0, 0, 0, 5, 99, 104, 111, 107, 101, 2, 0, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 37, 3, 0, 0, 0, 8, 63, 214, 99, 31, 138, 9, 2, 222}
		decodeString(&bs)
		decodeInteger(&bs)
		decodeFloat(&bs)
	}
}
