package repository

import (
	"context"
	"time"

	"shorty/internal/models"
	"shorty/internal/payload"
)

type LinkRepo interface {
	CreateLink(ctx context.Context, link *models.Link) (*models.Link, error)
	GetLinks(ctx context.Context, limit, offset int) ([]models.Link, error)
	GetLinkHash(ctx context.Context, hash string) (*models.Link, error)
	UpdateLink(ctx context.Context, link *models.Link) (*models.Link, error)
	DeleteLink(ctx context.Context, linkID uint) error
	CountLinks(ctx context.Context) (int64, error)
	FindLinkByID(ctx context.Context, linkID uint) (*models.Link, error)
	BlockLink(ctx context.Context, link *models.Link) (*models.Link, error)
	UnBlockLink(ctx context.Context, link *models.Link) (*models.Link, error)
	GetBlockedLinksCount(ctx context.Context) (int64, error)
	GetDeletedLinksCount(ctx context.Context) (int64, error)
	GetTotalLinks(ctx context.Context) (int64, error)
}

type StatRepo interface {
	AddClick(ctx context.Context, linkID uint) error
	GetClickedLinkStats(ctx context.Context, by string, from, to time.Time) []payload.GetStatsResponse
	GetAllLinksStats(ctx context.Context, from, to time.Time) []payload.LinkStatsResponse
}

type UserRepo interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUsers(ctx context.Context) ([]*models.User, error)
	GetUserByID(ctx context.Context, userID uint) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, userID uint) error
	BlockUsers(ctx context.Context, user *models.User) (*models.User, error)
	GetBlockedUsersCount(ctx context.Context) (int64, error)
	UnBlockUsers(ctx context.Context, user *models.User) (*models.User, error)
}
