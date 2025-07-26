package service

import (
	"context"
	"shorty/internal/models"
	"shorty/internal/payload"
	"time"
)

type LinkServ interface {
	Create(ctx context.Context, link *models.Link) (*models.Link, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.Link, error)
	GetByHash(ctx context.Context, hash string) (*models.Link, error)
	Update(ctx context.Context, link *models.Link) (*models.Link, error)
	Delete(ctx context.Context, linkID uint) error
	Count(ctx context.Context) (int64, error)
	FindByID(ctx context.Context, linkID uint) (*models.Link, error)
	Block(ctx context.Context, linkID uint) (*models.Link, error)
	UnBlock(ctx context.Context, linkID uint) (*models.Link, error)
	GetBlockedLinksCount(ctx context.Context) (int64, error)
	GetDeletedLinksCount(ctx context.Context) (int64, error)
	GetTotalLinks(ctx context.Context) (int64, error)
}

type StatServ interface {
	AddClick(ctx context.Context)
	GetClickedLinkStats(ctx context.Context, by string, from, to time.Time) []payload.GetStatsResponse
	GetAllLinksStats(ctx context.Context, from, to time.Time) []payload.LinkStatsResponse
}

type UserServ interface {
	GetAll(ctx context.Context, limit, offset int) ([]*models.User, error)
	GetByID(ctx context.Context, userID uint) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID uint) error
	Block(ctx context.Context, userID uint) (*models.User, error)
	UnBlock(ctx context.Context, userID uint) (*models.User, error)
	GetBlockedUsersCount(ctx context.Context) (int64, error)
}

type AuthServ interface {
	Registration(ctx context.Context, name, email, password string, role models.Role, isBlocked bool) (*models.User, error)
	Login(ctx context.Context, email, password string, role models.Role) (*models.User, error)
}
