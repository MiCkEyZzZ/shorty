package payload

import "shorty/internal/models"

// CreateLinkRequest represents the request payload for creating a new shortened link.
type CreateLinkRequest struct {
	URL string `json:"url" validate:"required,url"`
}

// UpdateLinkRequest represents the request payload for updating an existing shortened link.
type UpdateLinkRequest struct {
	URL       string `json:"url" validate:"required,url"`
	Hash      string `json:"hash"`
	IsBlocked bool   `json:"is_blocked"`
}

// BlockLinkRequest represents the request payload for blocking or unblocking a shortened link.
type BlockLinkRequest struct {
	IsBlocked bool `json:"is_blocked"`
}

// GetAllLinksResponse represents the response payload containing a list of links and their count.
type GetAllLinksResponse struct {
	Count int64         `json:"count"`
	Links []models.Link `json:"url"`
}
