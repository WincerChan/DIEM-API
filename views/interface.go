package views

import "github.com/gin-gonic/gin"

func Register(r *gin.Engine) {
	r.GET("/hitokoto/v2/", Hitokoto)
	r.GET("/gaviews/v1/", GAViews)
	r.GET("/blog-search/v1/", BlogSearchViews)
}
