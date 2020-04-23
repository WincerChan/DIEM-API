package gaviews

type Params struct {
	Prefix string `form:"prefix"`
}

type View struct {
	Path  string `json:"path"`
	Count int    `json:"count"`
}
