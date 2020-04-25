package gaviews

import (
	"testing"
)

func TestGetPathReport(t *testing.T) {
	var (
		in = &Params{Prefix: "/posts/"}
	)
	report := GetPathReport(in, "153425181")
	if report.DimensionFilterClauses[0].Filters[0].Expressions[0] != "/posts/" {
		t.Errorf("Error get report.")
	}
}
