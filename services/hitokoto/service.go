package hitokoto

import (
	C "DIEM-API/config"
	"encoding/json"
	"errors"
	"math/rand"
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

// override Scan implementation for Row.
func (h HitoInfo) Value() ([]byte, error) {
	result, err := json.Marshal(h)
	return result, err
}

// override Scan implementation for Row.
func (h *HitoInfo) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &h)
}

// fetch hitokoto from database
func FetchHitokoto(info *HitoInfo, length int) {
	seed := rand.Float32()
	row := C.PGConn.QueryRow("SELECT RANDOMFETCH($1, $2);", length, seed)
	row.Scan(info)
}
