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
const MAXKEYWORDLENGTH = 37

type Params struct {
	Paginate  []int    `form:"pages" json:"pages" deserialize:"BindPage"`
	Terms     []string `form:"terms" json:"terms" deserialize:"BindTerms"`
	Query     []string `form:"q" json:"q" deserialize:"BindQ"`
	DateRange []int    `form:"range" json:"range" deserialize:"BindRange"`
}

func (p *Params) BindPage(str string) string {
	pages := strings.Split(str, "-")
	p.Paginate = make([]int, 2)
	if len(pages) != 2 {
		return "invalid pages format, expected likes: 1-10"
	} else if len(pages[1]) > 1 {
		return "invalid pages format, page cannot greate than 10"
	}
	for i, pag := range pages {
		n, err := strconv.ParseUint(pag, 10, 32)
		if err != nil {
			return err.Error()
		}
		p.Paginate[i] = int(n)
	}
	return ""
}

func (p *Params) BindRange(str string) string {
	ranges := strings.Split(str, "~")
	if len(ranges) == 1 {
		return ""
	}
	p.DateRange = []int{0, int(time.Now().Unix())}
	for i, r := range ranges {
		if r != "" {
			t, err := time.Parse("2006-01-02", r)
			if err != nil {
				return err.Error()
			}
			p.DateRange[i] = int(t.Unix())
		}
	}
	return ""
}

func (p *Params) BindTerms(str string) string {
	terms := strings.Split(str, " ")
	p.Terms = make([]string, 0, 4)
	truncated := make([]string, 0, 4)
	if terms[0] == "" {
		return ""
	} else if len(terms) > 4 {
		truncated = terms[:4]
	} else {
		truncated = terms
	}
	for _, term := range truncated {
		if !strings.HasPrefix(term, "tags:") && !strings.HasPrefix(term, "category:") {
			return "invalid terms, expects category or tags"
		}
		p.Terms = append(p.Terms, term)
	}
	return ""
}

func (p *Params) BindQ(str string) string {
	q := []rune(str)
	if len(q) >= MAXKEYWORDLENGTH {
		p.Query = strings.Split(string(q[:MAXKEYWORDLENGTH]), ",")
	} else {
		p.Query = strings.Split(str, ",")
	}
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
