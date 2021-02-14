package views

import (
	M "DIEM-API/models"
	G "DIEM-API/models/googleanalytics"

	"github.com/gin-gonic/gin"
)

func checkGAParams(ctx *gin.Context, p *G.Params) {
	err := ctx.Bind(p)

	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		ctx.Abort()
	}
}

func GAViews(ctx *gin.Context) {
	p := new(G.Params)
	checkGAParams(ctx, p)
	pageView := M.GetReport(p)
	ctx.JSON(200, pageView)
}
