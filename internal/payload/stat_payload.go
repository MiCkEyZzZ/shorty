package payload

// GetStatsResponse represents a response containing general statistics over a specific period.
type GetStatsResponse struct {
	Period string `json:"period"`
	Sum    string `json:"sum"`
}

// LinkStatsResponse represents detailed statistics for a specific shortened link.
type LinkStatsResponse struct {
	LinkID        uint   `json:"link_id"`
	URL           string `json:"url"`
	TotalClicks   int64  `json:"total_clicks"`
	LastClickDate string `json:"last_click_date"`
	BlockedCount  int64  `json:"blocked_count"`
}
