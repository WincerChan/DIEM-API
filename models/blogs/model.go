package blogs

import (
	T "DIEM-API/tools"
	"log"
	"strings"
	"time"

	meili "github.com/meilisearch/meilisearch-go"
)

const INDEXUID = "blogs"

type Params struct {
	Paginate  string `form:"page"`
	Terms     string `form:"terms"`
	Query     string `form:"query"`
	DateRange string `form:"range"`
}

type Blog struct {
	Title    string   `json:"title"`
	Category string   `json:"category"`
	Date     string   `json:"date" time_format:"2006-01-02"`
	Tags     []string `json:"tags"`
	Snippet  string   `json:"snippet"`
	Url      string   `json:"url"`
}

type QueryRequest meili.SearchRequest

func NewQuery(q string) QueryRequest {
	r := meili.SearchRequest{
		AttributesToHighlight: []string{"*"},
		Matches:               true,
	}
	if q != "" {
		r.Query = q
	}
	return QueryRequest(r)
}

func (q *QueryRequest) AddTermsCondition(terms []string) {
	q.FacetFilters = terms
}

func (q *QueryRequest) AddDateFilter(start, end int64) {
	var lt, gt string
	if start < 0 {
		start = 0
	}
	gt = strings.Join([]string{"date", T.Str(start)}, ">=")
	if end <= 0 {
		end = time.Now().Unix()
	}
	lt = strings.Join([]string{"date", T.Str(end)}, "<")
	q.Filters = strings.Join([]string{lt, gt}, " AND ")
}

func (q *QueryRequest) AddPaginator(page, size int64) {
	if size > 0 {
		q.Limit = size
	} else {
		q.Limit = 10
	}
	if page > 1 {
		q.Offset = (page - 1) * size
	}
}

func genSnippet(text string, indeces []interface{}) string {
	snippets := make([]string, 0, 3)
	source, end := []rune(text), 0
	preStart := 0
	for i, index := range indeces {
		content := index.(map[string]interface{})
		start := int(content["start"].(float64)) + i*9
		log.Println(start, end)
		if start < preStart {
			continue
		}
		leftPart := []rune(text[:start])
		realStart := T.Max(0, len(leftPart)-50)
		end = T.Min(len(leftPart)+60, len(source)-1)
		if len(snippets) < 2 {
			snippets = append(snippets, string(source[realStart:end]))
			preStart = start + end
		}
	}
	snip := strings.Join(snippets, "...")
	if len(snip) < 360 {
		snip = strings.Join([]string{snip, string(source[end : end+160-len(snip)/3]), "..."}, "")
	} else {
		snip = strings.Join([]string{snip, "..."}, "")
	}
	return snip
}

func FormatHitsToBlog(hits []interface{}) []Blog {
	r := make([]Blog, 0, len(hits))
	for _, hit := range hits {
		b := new(Blog)
		response := hit.(map[string]interface{})
		formatted := response["_formatted"].(map[string]interface{})
		matches := response["_matchesInfo"].(map[string]interface{})
		if matches["content"] == nil {
			rawStr := []rune(formatted["content"].(string))
			b.Snippet = strings.Join([]string{string(rawStr[:T.Min(160, len(rawStr))]), "..."}, "")
		} else {
			b.Snippet = genSnippet(formatted["content"].(string), matches["content"].([]interface{}))
		}
		for _, tag := range formatted["tags"].([]interface{}) {
			b.Tags = append(b.Tags, tag.(string))
		}
		b.Category, _ = formatted["category"].(string)
		b.Title = formatted["title"].(string)
		b.Url = formatted["url"].(string)
		t := time.Unix(int64(formatted["date"].(float64)), 0)
		b.Date = t.Format("2006-01-02")
		r = append(r, *b)
	}
	return r
}
