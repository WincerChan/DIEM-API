package views

import (
	C "DIEM-API/config"
	service "DIEM-API/services/hitokoto"
	"bytes"

	"github.com/gin-gonic/gin"
)

func JSONFormat(ctx *gin.Context, info *C.HitoInfo) {
	ctx.JSON(200, info)
}

func PlainFormat(ctx *gin.Context, info *C.HitoInfo) {
	ctx.String(200, info.Hito+"——「"+info.Source+"」")
}

func JSONP(ctx *gin.Context, info *C.HitoInfo) {
	ctx.JSONP(200, info)
}

func JSFormat(ctx *gin.Context, info *C.HitoInfo) {
	var buf bytes.Buffer
	buf.WriteString("var hitokoto=\"")
	buf.WriteString(info.Hito)
	buf.WriteString("——「")
	buf.WriteString(info.Source)
	buf.WriteString("」\";var dom=document.querySelector('.hitokoto');")
	buf.WriteString("Array.isArray(dom)?dom[0].innerText=hitokoto:dom.innerText=hitokoto;")
	ctx.Data(200, "text/javascript; charset=utf-8", buf.Bytes())
}

// attempt to bind url params
func checkParams(ctx *gin.Context, p *C.Params) {
	err := ctx.Bind(p)

	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		ctx.Abort()
	}
}

func Hitokoto(ctx *gin.Context) {
	p := new(C.Params)

	checkParams(ctx, p)
	info := service.FetchHitokoto(p.Length)
	if p.Callback != "" {
		JSONP(ctx, info)
	} else if p.Encode == "js" {
		JSFormat(ctx, info)
	} else if p.Encode == "json" {
		JSONFormat(ctx, info)
	} else {
		PlainFormat(ctx, info)
	}
}
