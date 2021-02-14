package views

import (
	M "DIEM-API/models"
	H "DIEM-API/models/hitokoto"
	T "DIEM-API/tools"
	"bytes"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	bolt "go.etcd.io/bbolt"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func JSONFormat(ctx *gin.Context, info H.HitoInfo) {
	ctx.JSON(200, info)
}

func PlainFormat(ctx *gin.Context, info H.HitoInfo) {
	ctx.String(200, info.Hito+"——「"+info.Source+"」")
}

func JSONP(ctx *gin.Context, info H.HitoInfo) {
	ctx.JSONP(200, info)
}

func JSFormat(ctx *gin.Context, info H.HitoInfo) {
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
func checkParams(ctx *gin.Context, p *H.Params) {
	err := ctx.Bind(p)

	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		ctx.Abort()
	}
}

func fetchHitokoto(length int) (record H.HitoInfo) {
	randomNumber := r.Intn(H.IndexOf(length))
	key := T.Int32ToBytes(randomNumber)
	M.BoltDB.Read(func(tx *bolt.Tx) error {
		b := tx.Bucket(H.HitoBucket)
		record = H.LoadRecordFromBytes(b.Get(key)).Hitokoto
		return nil
	})
	return record
}

func Hitokoto(ctx *gin.Context) {
	p := new(H.Params)

	checkParams(ctx, p)
	info := fetchHitokoto(p.Length)
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
