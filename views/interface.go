package views

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, service string) {
	if strings.HasPrefix("hitokoto", service) {
		r.GET("/hitokoto/v2/", Hitokoto)
	}
	if strings.HasPrefix("analytics", service) {
		r.GET("/gaviews/v1/", GAViews)
	}
	if strings.HasPrefix("search", service) {
		r.GET("/blog-search/v1/", BlogSearchViews)
	}
}
