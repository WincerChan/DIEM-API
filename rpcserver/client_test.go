package rpcserver

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	times, _ := strconv.Atoi(os.Args[1])
	start := time.Now()
	var b []byte
	for i := 0; i < times; i++ {
		s := strconv.Itoa(37)
		b = append(b, []byte(s)...)
		f := fmt.Sprintf("%f", 0.3497)
		b = append(b, []byte(f)...)
		i := strconv.Itoa(37)
		b = append(b, []byte(i)...)
		// bf := new(bytes.Buffer)
		// encodeString(bf, "choke")
		// encodeInteger(bf, 37)
		// encodeFloat(bf, 0.3498)
		b = []byte{}
	}
	log.Println(time.Since(start))
	log.Println(b)
}
