package api

type Query struct {
	Query      string           `json:"query"`
	Parameters []QueryParameter `json:"parameters"`
}

type QueryParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
