package payload

type GetStatsResponse struct {
	Period string `json:"period"`
	Sum    string `json:"sum"`
}
