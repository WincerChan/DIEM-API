package models

import (
	G "DIEM-API/models/googleanalytics"
	T "DIEM-API/tools"
	"context"

	gar "google.golang.org/api/analyticsreporting/v4"
	"google.golang.org/api/option"
)

var (
	AnalyticsReportingService *gar.Service
)

func InitGACredential(path string) {
	ctx := context.Background()
	json := T.LoadJSON(path)
	ars, err := gar.NewService(ctx, option.WithCredentialsJSON(json))
	AnalyticsReportingService = ars
	T.CheckFatalError(err, false)
}

func GetReport(p *G.Params) G.ReportResponse {
	report := G.ConstructReport(p)
	req := &gar.GetReportsRequest{
		ReportRequests: []*gar.ReportRequest{
			report,
		},
	}
	resp, err := AnalyticsReportingService.Reports.BatchGet(req).Do()
	T.CheckException(err, "Failed to get analytics report.")
	return G.SimplifiedResponse(resp)
}
