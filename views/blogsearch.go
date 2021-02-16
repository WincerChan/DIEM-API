package views

import (
	M "DIEM-API/models"
	B "DIEM-API/models/blogs"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	meili "github.com/meilisearch/meilisearch-go"
)

func checkSearchParams(ctx *gin.Context, p *B.Params) {
	err := ctx.Bind(p)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		ctx.Abort()
	}
}

func validDateRange(r string) (begin, end int64) {
	f := func(s *int64, t time.Time, err error) {
		if err == nil {
			*s = t.Unix()
		}
	}
	lenOfTime := 8
	if len(r) <= lenOfTime {
		return
	}
	if r[lenOfTime] == '-' {
		beginTime, err := time.Parse("20060102", r[:lenOfTime])
		f(&begin, beginTime, err)
		r = r[lenOfTime:]
		return
	}
	if r[0] == '-' {
		endTime, err := time.Parse("20060102", r[1:lenOfTime+1])
		log.Println(endTime)
		f(&end, endTime, err)
	}
	return
}

func constructQuery(p *B.Params) B.QueryRequest {
	r := B.NewQuery(p.Query)
	if terms := strings.Split(p.Terms, " "); p.Terms != "" && len(terms) != 0 {
		r.AddTermsCondition(terms)
	}
	if paginate := strings.Split(p.Paginate, ":"); len(paginate) == 2 {
		page, _ := strconv.Atoi(paginate[0])
		size, _ := strconv.Atoi(paginate[1])
		r.AddPaginator(int64(page), int64(size))
	} else {
		r.AddPaginator(1, 10)
	}

	start, end := validDateRange(p.DateRange)
	r.AddDateFilter(start, end)
	return r
}

func fetchResponse(r B.QueryRequest) ([]B.Blog, int64) {
	req := meili.SearchRequest(r)
	resp := M.Execute(B.INDEXUID, meili.SearchRequest(req))
	return B.FormatHitsToBlog(resp.Hits), resp.NbHits
}

func BlogSearchViews(ctx *gin.Context) {
	p := new(B.Params)
	checkSearchParams(ctx, p)
	req := constructQuery(p)
	resp, n := fetchResponse(req)
	ctx.JSON(200, gin.H{
		"data":   resp,
		"counts": n,
	})
}
