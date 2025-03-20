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

// CreateLink добавляет новую ссылку в базу данных.
func (r *LinkRepository) CreateLink(ctx context.Context, link *models.Link) (*models.Link, error) {
	res := r.Database.DB.WithContext(ctx).Create(link)
	if res.Error != nil {
		logger.Error("Ошибка создания ссылки", zap.Error(res.Error))
		return nil, fmt.Errorf("ошибка при сохранении ссылки в БД: %w", res.Error)
	}
	logger.Info("Ссылка успешно создана", zap.Uint("linkID", link.ID))
	return link, nil
}

// GetLinks получает список ссылок с пагинацией.
func (r *LinkRepository) GetLinks(ctx context.Context, limit, offset int) ([]models.Link, error) {
	var links []models.Link
	res := r.Database.DB.
		Model(&models.Link{}).
		WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Find(&links)
	if res.Error != nil {
		logger.Error("Ошибка получения списка ссылок", zap.Error(res.Error))
		return nil, fmt.Errorf("ошибка при получении списка ссылок: %w", res.Error)
	}
	logger.Info("Получено ссылок", zap.Int("count", len(links)))
	return links, nil
}

// GetLinkHash ищет ссылку по хэшу.
func (r *LinkRepository) GetLinkHash(ctx context.Context, hash string) (*models.Link, error) {
	var link models.Link
	res := r.Database.DB.WithContext(ctx).First(&link, "hash = ?", hash)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Error("Ошибка поиска ссылки", zap.String("hash", hash), zap.Error(res.Error))
		return nil, res.Error
	}
	logger.Info("Ссылка найдена по хэшу", zap.String("hash", hash), zap.Uint("linkID", link.ID))
	return &link, nil
}

// UpdateLink обновляет данные ссылки в базе.
func (r *LinkRepository) UpdateLink(ctx context.Context, link *models.Link) (*models.Link, error) {
	res := r.Database.DB.WithContext(ctx).Clauses(clause.Returning{}).Updates(link)
	if res.Error != nil {
		logger.Error("Ошибка обновления ссылки", zap.Uint("linkID", link.ID), zap.Error(res.Error))
		return nil, fmt.Errorf("ошибка при обновлении ссылки в БД: %w", res.Error)
	}
	logger.Info("Ссылка успешно обновлена", zap.Uint("linkID", link.ID))
	return link, nil
}

// DeleteLink удаляет ссылку по идентификатору.
func (r *LinkRepository) DeleteLink(ctx context.Context, linkID uint) error {
	res := r.Database.DB.WithContext(ctx).Delete(&models.Link{}, linkID)
	if res.Error != nil {
		logger.Error("Ошибка удаления ссылки", zap.Uint("linkID", linkID), zap.Error(res.Error))
		return fmt.Errorf("ошибка при удалении ссылки из БД: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		logger.Warn("Ссылка не найдена для удаления", zap.Uint("linkID", linkID))
		return fmt.Errorf("ссылка с ID %d не найдена", linkID)
	}
	logger.Info("Ссылка успешно удалена", zap.Uint("linkID", linkID))
	return nil
}

// CountLink возвращает количество ссылок в базе.
func (r *LinkRepository) CountLink(ctx context.Context) (int64, error) {
	var count int64
	res := r.Database.DB.
		WithContext(ctx).
		Model(&models.Link{}).
		Where("deleted_at IS NULL").
		Count(&count)
	if res.Error != nil {
		logger.Error("Ошибка подсчёта ссылок", zap.Error(res.Error))
		return 0, fmt.Errorf("ошибка при подсчёте ссылок: %w", res.Error)
	}
	logger.Info("Количество ссылок", zap.Int64("count", count))
	return count, nil
}

// FindLinkByID ищет ссылку по идентификатору.
func (r *LinkRepository) FindLinkByID(ctx context.Context, linkID uint) (*models.Link, error) {
	var link models.Link
	res := r.Database.DB.WithContext(ctx).First(&link, linkID)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			logger.Warn("Ссылка не найдена", zap.Uint("linkID", linkID))
			return nil, nil
		}
		logger.Error("Ошибка при поиске ссылки", zap.Uint("linkID", linkID), zap.Error(res.Error))
		return nil, fmt.Errorf("ошибка при поиске ссылки: %w", res.Error)
	}
	logger.Info("Ссылка найдена", zap.Uint("linkID", linkID))
	return &link, nil
}

// BlockLink блокирует ссылку по идентификатору.
func (r *LinkRepository) BlockLink(ctx context.Context, link *models.Link) (*models.Link, error) {
	res := r.Database.DB.WithContext(ctx).
		Model(&models.Link{}).
		Where("id = ?", link.ID).
		Updates(map[string]interface{}{
			"is_blocked": link.IsBlocked,
		})

	if res.Error != nil {
		logger.Error("Ошибка блокировки ссылки", zap.Uint("id", link.ID), zap.Error(res.Error))
		return nil, fmt.Errorf("ошибка при блокировке ссылки: %w", res.Error)
	}

	logger.Info("Ссылка заблокирована", zap.Uint("id", link.ID), zap.Bool("is_blocked", link.IsBlocked))
	return link, nil
}

// UnBlockLink заблокирует ссылку по идентификатору.
func (r *LinkRepository) UnBlockLink(ctx context.Context, link *models.Link) (*models.Link, error) {
	res := r.Database.DB.WithContext(ctx).
		Model(&models.Link{}).
		Where("id = ?", link.ID).
		Updates(map[string]interface{}{
			"is_blocked": link.IsBlocked,
		})

	if res.Error != nil {
		logger.Error("Ошибка снятие блокировки с ссылки", zap.Uint("id", link.ID), zap.Error(res.Error))
		return nil, fmt.Errorf("ошибка при снятии блокировки c ссылки: %w", res.Error)
	}

	logger.Info("Ссылка разблокирована", zap.Uint("id", link.ID), zap.Bool("is_blocked", link.IsBlocked))
	return link, nil
}
