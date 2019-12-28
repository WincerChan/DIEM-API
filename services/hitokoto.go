package services

import (
	C "DIEM-API/config"
	"bytes"
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
)

var params struct {
	Length   int    `form:"length"`
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

func JSONFormat(ctx *gin.Context, info *HitoInfo) {
	ctx.JSON(200, info)
}

func JSFormat(ctx *gin.Context, info *HitoInfo) {
	var buf bytes.Buffer
	buf.WriteString("var hitokoto=\"")
	buf.WriteString(info.Hito)
	buf.WriteString("——「")
	buf.WriteString(info.Source)
	buf.WriteString("」\";var dom=document.querySelector('.hitokoto');")
	buf.WriteString("Array.isArray(dom)?dom[0].innerText=hitokoto:dom.innerText=hitokoto;")
	ctx.Data(200, "text/javascript; charset=utf-8", buf.Bytes())
}

func PlainFormat(ctx *gin.Context, info *HitoInfo) {
	ctx.String(200, info.Hito+"——「"+info.Source+"」")
}

func validateQueryParams(ctx *gin.Context) bool {
	err := ctx.Bind(&params)

	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return false
	}
	return true
}

func Hitokoto(ctx *gin.Context) {
	if !validateQueryParams(ctx) {
		return
	}
	info := new(HitoInfo)
	fetch(info)
	if params.Callback != "" {
		ctx.JSONP(200, info)
		return
	}

	switch params.Encode {
	case "json":
		JSONFormat(ctx, info)
	case "js":
		JSFormat(ctx, info)
	default:
		PlainFormat(ctx, info)
	}
}
