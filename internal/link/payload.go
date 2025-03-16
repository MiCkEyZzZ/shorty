package link

// LinkCreateRequest - структура для создания новой сокращенной ссылки.
type LinkCreateRequest struct {
	URL string `json:"url" validate:"required,url"`
}

// LinkUpdateRequest - структура для обновления сокращенной ссылки.
type LinkUpdateRequest struct {
	URL  string `json:"url" validate:"required,url"`
	Hash string `json:"hash"`
}
