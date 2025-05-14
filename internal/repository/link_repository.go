package repository

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"shorty/internal/models"
	"shorty/pkg/db"
	"shorty/pkg/logger"
)

// LinkRepository handles database operations for the Link entity.
type LinkRepository struct {
	Database *db.DB
}

// NewLinkRepository creates and returns a new instance of LinkRepository.
func NewLinkRepository(db *db.DB) *LinkRepository {
	return &LinkRepository{Database: db}
}

// CreateLink creates a new shortened link record in the database.
func (r *LinkRepository) CreateLink(ctx context.Context, link *models.Link) (*models.Link, error) {
	result := r.Database.DB.WithContext(ctx).Create(link)
	if result.Error != nil {
		logger.Error("Failed to create link", zap.Error(result.Error))
		return nil, fmt.Errorf("failed to save link in the database: %w", result.Error)
	}
	logger.Info("Link successfully created", zap.Uint("linkID", link.ID))
	return link, nil
}

// GetLinks retrieves a paginated list of active, unblocked links.
func (r *LinkRepository) GetLinks(ctx context.Context, limit, offset int) ([]models.Link, error) {
	var links []models.Link
	result := r.Database.DB.
		Model(&models.Link{}).
		WithContext(ctx).
		Where("deleted_at IS NULL AND is_blocked = false").
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Find(&links)
	if result.Error != nil {
		logger.Error("Failed to retrieve link list", zap.Error(result.Error))
		return nil, fmt.Errorf("failed to retrieve link list: %w", result.Error)
	}
	logger.Info("Links retrieved", zap.Int("count", len(links)))
	return links, nil
}

// GetLinkHash retrieves a link by its hash if it is not blocked.
func (r *LinkRepository) GetLinkHash(ctx context.Context, hash string) (*models.Link, error) {
	var link models.Link
	result := r.Database.DB.WithContext(ctx).
		Where("hash = ? AND is_blocked = false", hash).
		First(&link)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		logger.Error("Failed to find link by hash", zap.String("hash", hash), zap.Error(result.Error))
		return nil, result.Error
	}
	logger.Info("Link found by hash", zap.String("hash", hash), zap.Uint("linkID", link.ID))
	return &link, nil
}

// UpdateLink updates an existing link and returns the updated record.
func (r *LinkRepository) UpdateLink(ctx context.Context, link *models.Link) (*models.Link, error) {
	result := r.Database.DB.WithContext(ctx).Clauses(clause.Returning{}).Updates(link)
	if result.Error != nil {
		logger.Error("Failed to update link", zap.Uint("linkID", link.ID), zap.Error(result.Error))
		return nil, fmt.Errorf("failed to update link in the database: %w", result.Error)
	}
	logger.Info("Link successfully updated", zap.Uint("linkID", link.ID))
	return link, nil
}

// DeleteLink marks a link as deleted by setting the deleted_at timestamp.
func (r *LinkRepository) DeleteLink(ctx context.Context, linkID uint) error {
	result := r.Database.DB.WithContext(ctx).
		Model(&models.Link{}).
		Where("id = ?", linkID).
		Update("deleted_at", gorm.Expr("Now()"))
	if result.Error != nil {
		logger.Error("Failed to delete link", zap.Uint("linkID", linkID), zap.Error(result.Error))
		return fmt.Errorf("failed to delete link from database: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		logger.Warn("Link not found for deletion", zap.Uint("linkID", linkID))
		return fmt.Errorf("link with ID %d not found", linkID)
	}
	logger.Info("Link successfully deleted", zap.Uint("linkID", linkID))
	return nil
}

// CountLinks returns the number of non-deleted links.
func (r *LinkRepository) CountLinks(ctx context.Context) (int64, error) {
	var count int64
	result := r.Database.DB.
		WithContext(ctx).
		Model(&models.Link{}).
		Where("deleted_at IS NULL").
		Count(&count)
	if result.Error != nil {
		logger.Error("Failed to count links", zap.Error(result.Error))
		return 0, fmt.Errorf("failed to count links: %w", result.Error)
	}
	logger.Info("Total active links count", zap.Int64("count", count))
	return count, nil
}

// FindLinkByID finds a link by its unique ID.
func (r *LinkRepository) FindLinkByID(ctx context.Context, linkID uint) (*models.Link, error) {
	var link models.Link
	result := r.Database.DB.WithContext(ctx).First(&link, linkID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Warn("Link not found", zap.Uint("linkID", linkID))
			return nil, nil
		}
		logger.Error("Error retrieving link", zap.Uint("linkID", linkID), zap.Error(result.Error))
		return nil, fmt.Errorf("error retrieving link: %w", result.Error)
	}
	logger.Info("Link found", zap.Uint("linkID", linkID))
	return &link, nil
}

// BlockLink sets the 'is_blocked' flag to true for a link.
func (r *LinkRepository) BlockLink(ctx context.Context, link *models.Link) (*models.Link, error) {
	res := r.Database.DB.WithContext(ctx).
		Model(&models.Link{}).
		Where("id = ?", link.ID).
		Updates(map[string]interface{}{"is_blocked": true})

	if res.Error != nil {
		logger.Error("Failed to block link", zap.Uint("id", link.ID), zap.Error(res.Error))
		return nil, fmt.Errorf("failed to block link: %w", res.Error)
	}

	logger.Info("Link blocked", zap.Uint("id", link.ID), zap.Bool("is_blocked", true))
	return link, nil
}

// UnBlockLink sets the 'is_blocked' flag to false for a link.
func (r *LinkRepository) UnBlockLink(ctx context.Context, link *models.Link) (*models.Link, error) {
	res := r.Database.DB.WithContext(ctx).
		Model(&models.Link{}).
		Where("id = ?", link.ID).
		Updates(map[string]interface{}{"is_blocked": false})

	if res.Error != nil {
		logger.Error("Failed to unblock link", zap.Uint("id", link.ID), zap.Error(res.Error))
		return nil, fmt.Errorf("failed to unblock link: %w", res.Error)
	}

	logger.Info("Link unblocked", zap.Uint("id", link.ID), zap.Bool("is_blocked", false))
	return link, nil
}

// GetBlockedLinksCount returns the number of blocked links.
func (r *LinkRepository) GetBlockedLinksCount(ctx context.Context) (int64, error) {
	var count int64
	result := r.Database.DB.WithContext(ctx).
		Model(&models.Link{}).
		Where("is_blocked = ?", true).
		Count(&count)
	if result.Error != nil {
		logger.Error("Failed to count blocked links", zap.Error(result.Error))
		return 0, fmt.Errorf("failed to count blocked links: %w", result.Error)
	}
	return count, nil
}

// GetDeletedLinksCount returns the number of deleted links.
func (r *LinkRepository) GetDeletedLinksCount(ctx context.Context) (int64, error) {
	var count int64
	result := r.Database.DB.WithContext(ctx).
		Model(&models.Link{}).
		Where("deleted_at IS NOT NULL").
		Count(&count)
	if result.Error != nil {
		logger.Error("Failed to count deleted links", zap.Error(result.Error))
		return 0, result.Error
	}
	return count, nil
}

// GetTotalLinks returns the total number of links ever created (including deleted).
func (r *LinkRepository) GetTotalLinks(ctx context.Context) (int64, error) {
	var count int64
	result := r.Database.DB.WithContext(ctx).Model(&models.Link{}).Unscoped().Count(&count)
	if result.Error != nil {
		logger.Error("Failed to count total links", zap.Error(result.Error))
		return 0, result.Error
	}
	return count, nil
}
