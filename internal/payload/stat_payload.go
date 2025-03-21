package payload

type GetStatsResponse struct {
	Period string `json:"period"`
	Sum    string `json:"sum"`
}

type LinkStatsResponse struct {
	LinkID        uint   `json:"link_id"`
	URL           string `json:"url"`
	TotalClicks   int64  `json:"total_clicks"`
	LastClickDate string `json:"last_click_date"`
	BlockedCount  int64  `json:"blocked_count"`
}
