package blogs

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const INDEXUID = "blogs"

type Params struct {
	Paginate  string `form:"pages" json:"pages" deserialize:"BindPage"`
	Terms     string `form:"terms" json:"terms" deserialize:"BindTerms"`
	Query     string `form:"q" json:"q" deserialize:"BindQ"`
	DateRange string `form:"range" json:"range" deserialize:"BindRange"`
}

func (p *Params) BindPage(str string) string {
	pages := strings.Split(str, "-")
	if len(pages) != 2 {
		return "invalid pages format, expected likes: 1-10"
	}
	for _, p := range pages {
		_, err := strconv.Atoi(p)
		if err != nil {
			return err.Error()
		}
	}
	p.Paginate = str
	return ""
}

func (p *Params) BindRange(str string) string {
	ranges := strings.Split(str, "~")
	times := []string{strconv.Itoa(0), strconv.Itoa(int(time.Now().Unix()))}
	if str == "" {
	} else if len(ranges) != 2 {
		return "invalid range format, expected likes: 2020-02-12~2021-03-24"
	}
	for i, r := range ranges {
		if r != "" {
			t, err := time.Parse("2006-01-02", r)
			if err != nil {
				return err.Error()
			}
			times[i] = strconv.Itoa(int(t.Unix()))
		}
	}
	p.DateRange = strings.Join(times, "~")
	return ""
}

func (p *Params) BindTerms(str string) string {
	if str == "" {
		return ""
	}
	terms := strings.Split(str, " ")
	if len(terms) > 5 {
		return "max allowed terms is 5"
	}
	for _, term := range terms {
		if !strings.HasPrefix(term, "tags:") && !strings.HasPrefix(term, "category:") {
			return "invalid terms, expects category or tags"
		}
	}
	p.Terms = str
	return ""
}

func (p *Params) BindQ(str string) string {
	p.Query = str
	return ""
}

func (p *Params) Serialize() string {
	paramsValue := reflect.ValueOf(*p)
	lens := paramsValue.NumField()
	results := make([]string, lens, lens)
	for i := 0; i < lens; i++ {
		results[i] = paramsValue.Field(i).String()
	}
	return strings.Join(results, "\x00")
}

func BindStruct(m url.Values, p *Params) error {
	paramKey := reflect.TypeOf(p).Elem()
	paramMethod := reflect.ValueOf(p)
	paramValue := paramMethod.Elem()
	for i := 0; i < paramValue.NumField(); i++ {
		field := paramKey.Field(i)
		key := field.Tag.Get("form")
		method := paramMethod.MethodByName(field.Tag.Get("deserialize"))
		err := method.Call([]reflect.Value{reflect.ValueOf(m.Get(key))})
		if !err[0].IsZero() {
			return errors.New(err[0].String())
		}
	}
	return nil
}
