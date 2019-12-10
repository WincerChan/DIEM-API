package services

import (
	C "DIEM-API/config"
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
)

var params struct {
	Length   int    `form:"length" binding:"required,min=1,max=100"`
	Callback string `form:"callback"`
	Encode   string `form:"encode"`
}

type HitoInfo struct {
	Source string `json:"source"`
	Hito   string `json:"hitokoto"`
}

func (h HitoInfo) Value() []byte {
	result, _ := json.Marshal(h)
	return result
}

func (h *HitoInfo) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &h)
}

func fetch(info *HitoInfo) {
	row := C.PGConn.QueryRow("SELECT RANDOMFETCH($1);", params.Length)
	row.Scan(info)
}

func Hitokoto(ctx *gin.Context) {
	info := new(HitoInfo)
	err := ctx.Bind(&params)

	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
	}
	fetch(info)
	ctx.JSON(200, info)
}
