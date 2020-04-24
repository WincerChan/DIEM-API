package views

import (
	service "DIEM-API/services/gaviews"
	"github.com/gin-gonic/gin"
)

func checkGAParams(ctx *gin.Context, p *service.Params) {
	err := ctx.Bind(p)

	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		ctx.Abort()
	}
}

func GAViews(ctx *gin.Context) {
	p := new(service.Params)
	pageView := new(service.ReportResponse)
	checkGAParams(ctx, p)

	service.GetPageViews(p, pageView)
	ctx.JSON(200, pageView)
}
