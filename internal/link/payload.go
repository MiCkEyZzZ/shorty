package link

// CreateLinkRequest - структура для создания новой сокращенной ссылки.
type CreateLinkRequest struct {
	URL string `json:"url" validate:"required,url"`
}

// UpdateLinkRequest - структура для обновления сокращенной ссылки.
type UpdateLinkRequest struct {
	URL  string `json:"url" validate:"required,url"`
	Hash string `json:"hash"`
}

type GetAllLinksResponse struct {
	Count int64  `json:"count"`
	Links []Link `json:"url"`
}
