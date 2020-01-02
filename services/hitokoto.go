package services

import (
	C "DIEM-API/config"
	"bytes"
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
)

type params struct {
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

func fetchHitokoto(info *HitoInfo, length int) {
	row := C.PGConn.QueryRow("SELECT RANDOMFETCH($1);", length)
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

func checkValidReq(ctx *gin.Context, p *params) {
	err := ctx.Bind(p)

	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		ctx.Abort()
	}
}

func Hitokoto(ctx *gin.Context) {
	p := new(params)
	info := new(HitoInfo)

	checkValidReq(ctx, p)
	fetchHitokoto(info, p.Length)

	if p.Callback != "" {
		ctx.JSONP(200, info)
	} else if p.Encode == "js" {
		JSFormat(ctx, info)
	} else if p.Encode == "json" {
		JSONFormat(ctx, info)
	} else {
		PlainFormat(ctx, info)
	}
}
