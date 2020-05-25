package gaviews

import (
	C "DIEM-API/config"
	T "DIEM-API/tools"
	gar "google.golang.org/api/analyticsreporting/v4"
)

// gaviews url params
type Params struct {
	Prefix string `form:"prefix"`
}

// gaviews response format
type ReportResponse struct {
	Details []accessInfo `json:"details"`
	Total   int          `json:"total"`
}

//gaviews every page access info format
type accessInfo struct {
	Path  string `json:"path"`
	Count int    `json:"count"`
}

// the google analytics request report template
func GetPathReport(p *Params, viewID string) *gar.ReportRequest {
	reportParams := new(gar.ReportRequest)
	reportParams.ViewId = viewID
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

// remove some field that don't need
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

func GetPageViews(p *Params, pageView *ReportResponse) {
	report := GetPathReport(p, C.GAViewID)
	req := &gar.GetReportsRequest{
		ReportRequests: []*gar.ReportRequest{
			report,
		},
	}
	resp, err := C.AnalyticsReportingService.Reports.BatchGet(req).Do()
	T.CheckException(err, "Failed to get analytics report.")
	*pageView = SimplifiedResponse(resp)
}
