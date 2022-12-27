package helpers

type ZSRequest struct {
	Source     []string       `json:"_source"`
	Explain    bool           `json:"explain"`
	From       int            `json:"from"`
	MaxResults int            `json:"max_results"`
	Query      ZSRequestQuery `json:"query"`
	SearchType string         `json:"search_type"`
	SortFields []string       `json:"sort_fields"`
}

type ZSRequestQuery struct {
	Term  string   `json:"term"`
	Terms []string `json:"terms"`
}
