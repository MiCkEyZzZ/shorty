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

// LinkRepository отвечает за операции с базой данных для сущности Link.
type LinkRepository struct {
	Database *db.DB
}

// NewLinkRepository создаёт новый экземпляр LinkRepository.
func NewLinkRepository(db *db.DB) *LinkRepository {
	return &LinkRepository{Database: db}
}

// CreateLink метод для создания новой ссылки.
func (r *LinkRepository) CreateLink(ctx context.Context, link *models.Link) (*models.Link, error) {
	result := r.Database.DB.WithContext(ctx).Create(link)
	if result.Error != nil {
		logger.Error("Ошибка создания ссылки", zap.Error(result.Error))
		return nil, fmt.Errorf("ошибка при сохранении ссылки в БД: %w", result.Error)
	}
	logger.Info("Ссылка успешно создана", zap.Uint("linkID", link.ID))
	return link, nil
}

// GetLinks метод для получения списка ссылок с пагинацией.
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
		logger.Error("Ошибка получения списка ссылок", zap.Error(result.Error))
		return nil, fmt.Errorf("ошибка при получении списка ссылок: %w", result.Error)
	}
	logger.Info("Получено ссылок", zap.Int("count", len(links)))
	return links, nil
}

// GetLinkHash метод для поиска ссылки по хэшу.
func (r *LinkRepository) GetLinkHash(ctx context.Context, hash string) (*models.Link, error) {
	var link models.Link
	result := r.Database.DB.WithContext(ctx).
		Where("hash = ? AND is_blocked = false", hash).
		First(&link)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		logger.Error("Ошибка поиска ссылки", zap.String("hash", hash), zap.Error(result.Error))
		return nil, result.Error
	}
	logger.Info("Ссылка найдена по хэшу", zap.String("hash", hash), zap.Uint("linkID", link.ID))
	return &link, nil
}

// UpdateLink метод для обновления ссылки.
func (r *LinkRepository) UpdateLink(ctx context.Context, link *models.Link) (*models.Link, error) {
	result := r.Database.DB.WithContext(ctx).Clauses(clause.Returning{}).Updates(link)
	if result.Error != nil {
		logger.Error("Ошибка обновления ссылки", zap.Uint("linkID", link.ID), zap.Error(result.Error))
		return nil, fmt.Errorf("ошибка при обновлении ссылки в БД: %w", result.Error)
	}
	logger.Info("Ссылка успешно обновлена", zap.Uint("linkID", link.ID))
	return link, nil
}

// DeleteLink метод для удаления ссылки по идентификатору.
func (r *LinkRepository) DeleteLink(ctx context.Context, linkID uint) error {
	result := r.Database.DB.WithContext(ctx).
		Model(&models.Link{}).
		Where("id = ?", linkID).
		Update("deleted_at", gorm.Expr("Now()"))
	if result.Error != nil {
		logger.Error("Ошибка удаления ссылки", zap.Uint("linkID", linkID), zap.Error(result.Error))
		return fmt.Errorf("ошибка при удалении ссылки из БД: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		logger.Warn("Ссылка не найдена для удаления", zap.Uint("linkID", linkID))
		return fmt.Errorf("ссылка с ID %d не найдена", linkID)
	}
	logger.Info("Ссылка успешно удалена", zap.Uint("linkID", linkID))
	return nil
}

// CountLinks метод для возврата количества ссылок.
func (r *LinkRepository) CountLinks(ctx context.Context) (int64, error) {
	var count int64
	result := r.Database.DB.
		WithContext(ctx).
		Model(&models.Link{}).
		Where("deleted_at IS NULL").
		Count(&count)
	if result.Error != nil {
		logger.Error("Ошибка подсчёта ссылок", zap.Error(result.Error))
		return 0, fmt.Errorf("ошибка при подсчёте ссылок: %w", result.Error)
	}
	logger.Info("Количество ссылок", zap.Int64("count", count))
	return count, nil
}

// FindLinkByID метод для поиска ссылок по идентификатору.
func (r *LinkRepository) FindLinkByID(ctx context.Context, linkID uint) (*models.Link, error) {
	var link models.Link
	result := r.Database.DB.WithContext(ctx).First(&link, linkID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Warn("Ссылка не найдена", zap.Uint("linkID", linkID))
			return nil, nil
		}
		logger.Error("Ошибка при поиске ссылки", zap.Uint("linkID", linkID), zap.Error(result.Error))
		return nil, fmt.Errorf("ошибка при поиске ссылки: %w", result.Error)
	}
	logger.Info("Ссылка найдена", zap.Uint("linkID", linkID))
	return &link, nil
}

// BlockLink метод для блокировки ссылки по идентификатору.
func (r *LinkRepository) BlockLink(ctx context.Context, link *models.Link) (*models.Link, error) {
	res := r.Database.DB.WithContext(ctx).
		Model(&models.Link{}).
		Where("id = ?", link.ID).
		Updates(map[string]interface{}{"is_blocked": true})

	if res.Error != nil {
		logger.Error("Ошибка блокировки ссылки", zap.Uint("id", link.ID), zap.Error(res.Error))
		return nil, fmt.Errorf("ошибка при блокировке ссылки: %w", res.Error)
	}

	logger.Info("Ссылка заблокирована", zap.Uint("id", link.ID), zap.Bool("is_blocked", true))
	return link, nil
}

// UnBlockLink метод для разблокировки ссылки по идентификатору.
func (r *LinkRepository) UnBlockLink(ctx context.Context, link *models.Link) (*models.Link, error) {
	res := r.Database.DB.WithContext(ctx).
		Model(&models.Link{}).
		Where("id = ?", link.ID).
		Updates(map[string]interface{}{"is_blocked": false})

	if res.Error != nil {
		logger.Error("Ошибка снятия блокировки с ссылки", zap.Uint("id", link.ID), zap.Error(res.Error))
		return nil, fmt.Errorf("ошибка при снятии блокировки с ссылки: %w", res.Error)
	}

	logger.Info("Ссылка разблокирована", zap.Uint("id", link.ID), zap.Bool("is_blocked", false))
	return link, nil
}

// GetBlockedLinksCount метод для получения количества заблокированных ссылок.
func (r *LinkRepository) GetBlockedLinksCount(ctx context.Context) (int64, error) {
	var count int64
	result := r.Database.DB.WithContext(ctx).
		Model(&models.Link{}).
		Where("is_blocked = ?", true).
		Count(&count)
	if result.Error != nil {
		logger.Error("Ошибка при получении количества заблокированных ссылок", zap.Error(result.Error))
		return 0, fmt.Errorf("ошибка при получении количества заблокированных ссылок: %w", result.Error)
	}
	return count, nil
}

// GetDeletedLinksCount метод для получения количества удалённых ссылок.
func (r *LinkRepository) GetDeletedLinksCount(ctx context.Context) (int64, error) {
	var count int64
	result := r.Database.DB.WithContext(ctx).
		Model(&models.Link{}).
		Where("deleted_at IS NOT NULL").
		Count(&count)
	if result.Error != nil {
		logger.Error("Ошибка при получении количества удаленных ссылок", zap.Error(result.Error))
		return 0, result.Error
	}
	return count, nil
}
