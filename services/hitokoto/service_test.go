package hitokoto

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"testing"
)

var PGConn *sqlx.DB

func init() {
	PGConn, _ = sqlx.Open("postgres", "port=5433")
}

func TestHitoInfo_Scan(t *testing.T) {
	var (
		in = &HitoInfo{
			Source: "今天天气真好",
			Hito:   "我也觉得",
		}
	)
	byte, err := json.Marshal(in)
	if err != nil {
		t.Errorf("Scan(%v) error; got %v", in, byte)
	}
	in.Scan(byte)
}

func TestHitoInfo_Value(t *testing.T) {
	var (
		in = &HitoInfo{
			Source: "今天天气真好",
			Hito:   "我也觉得",
		}
	)
	value, err := in.Value()
	if err != nil {
		t.Errorf("Value(%v) error; got %v", in, value)
	}
}

func TestFetchHitokoto(t *testing.T) {
	var (
		in = new(HitoInfo)
	)
	row := PGConn.QueryRow("SELECT RANDOMFETCH(1);")
	row.Scan(in)
}
