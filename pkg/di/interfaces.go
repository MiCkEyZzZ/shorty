package di

import "context"

type IStatService interface {
	AddClick(ctx context.Context, linkID uint) error
}
