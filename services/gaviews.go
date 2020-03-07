package services

import (
	C "DIEM-API/config"
	T "DIEM-API/tools"
	"context"
	"github.com/gin-gonic/gin"
	gar "google.golang.org/api/analyticsreporting/v4"
	"google.golang.org/api/option"
	"io/ioutil"
	"os"
)

type gaparams struct {
	Prefix string `form:"prefix"`
}

type View struct {
	Path  string `json:"path"`
	Count int    `json:"count"`
}

type ReportResponse struct {
	Details []View `json:"details"`
	Total   int    `json:"total"`
}

var analyticsreportingService *gar.Service

func init() {
	ctx := context.Background()
	json := parseCredentialJSON()
	analyticsreportingService, _ = gar.NewService(ctx, option.WithCredentialsJSON(json))
}

func parseCredentialJSON() []byte {
	jsonFile, _ := os.Open("./credential.json")

	byteValue, _ := ioutil.ReadAll(jsonFile)

	return byteValue
}

func getViews(ctx *gin.Context) *gar.ReportRequest {
	p := new(gaparams)
	err := ctx.Bind(p)

	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		ctx.Abort()
	}
	reportParams := new(gar.ReportRequest)
	reportParams.ViewId = C.GAViewID
	reportParams.DateRanges = []*gar.DateRange{&gar.DateRange{StartDate: "2017-06-18", EndDate: "today"}}
	reportParams.Metrics = []*gar.Metric{&gar.Metric{Expression: "ga:pageviews"}}
	reportParams.Dimensions = []*gar.Dimension{&gar.Dimension{Name: "ga:pagePath"}}
	reportParams.DimensionFilterClauses = []*gar.DimensionFilterClause{
		{
			Filters: []*gar.DimensionFilter{&gar.DimensionFilter{
				DimensionName: "ga:pagePath",
				Operator:      "BEGINS_WITH",
				Expressions:   []string{p.Prefix},
			},
			},
		},
	}

	return reportParams
}

func validateResponse(response *gar.GetReportsResponse) (rr ReportResponse) {
	for _, report := range response.Reports {
		rr.Total = T.Int(report.Data.Totals[0].Values[0])
		rr.Details = make([]View, report.Data.RowCount)
		for i, row := range report.Data.Rows {
			rr.Details[i] = View{Path: row.Dimensions[0], Count: T.Int(row.Metrics[0].Values[0])}
		}
	}
	return
}

func GAViews(ctx *gin.Context) {
	report := getViews(ctx)
	req := &gar.GetReportsRequest{
		ReportRequests: []*gar.ReportRequest{
			report,
		},
	}
	resp, err := analyticsreportingService.Reports.BatchGet(req).Do()
	T.CheckException(err, "Failed to get analytics report.")
	ctx.JSON(200, validateResponse(resp))
}
