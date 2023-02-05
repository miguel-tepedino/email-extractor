package types

type Query struct {
	SearchType string `json:"search_type"`
	From       int    `json:"from"`
	MaxResults int    `json:"max_results"`
	Source     []any  `json:"_source"`
	Search     any    `json:"query"`
}

type Search struct {
	Term   string `json:"term"`
	Offset uint   `json:"offset"`
}
