package googleanalytics

import (
	T "DIEM-API/tools"

	gar "google.golang.org/api/analyticsreporting/v4"
)

var (
	GAViewID string
)

type Params struct {
	Prefix string `form:"prefix"`
}

type ReportResponse struct {
	Details []accessInfo `json:"details"`
	Total   int          `json:"total"`
}

type accessInfo struct {
	Path  string `json:"path"`
	Count int    `json:"count"`
}

func ConstructReport(p *Params) *gar.ReportRequest {
	reportParams := new(gar.ReportRequest)
	reportParams.ViewId = GAViewID
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

func SimplifiedResponse(response *gar.GetReportsResponse) (rr ReportResponse) {
	rr = *new(ReportResponse)
	for _, report := range response.Reports {
		rr.Total = T.Int(report.Data.Totals[0].Values[0])
		rr.Details = make([]accessInfo, report.Data.RowCount)
		for i, row := range report.Data.Rows {
			rr.Details[i] = accessInfo{Path: row.Dimensions[0], Count: T.Int(row.Metrics[0].Values[0])}
		}
	}
	return
}
