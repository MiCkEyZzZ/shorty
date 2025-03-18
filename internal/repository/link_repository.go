package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"shorty/internal/models"
	"shorty/pkg/db"
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
		log.Printf("[LinkRepository] Ошибка создания ссылки: %v", res.Error)
		return nil, fmt.Errorf("ошибка при сохранении ссылки в БД: %w", res.Error)
	}
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
		log.Printf("[LinkRepository] Ошибка получения списка ссылок: %v", res.Error)
		return nil, fmt.Errorf("ошибка при получении списка ссылок: %w", res.Error)
	}
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
		log.Printf("[LinkRepository] Ошибка поиска ссылки: %v", res.Error)
		return nil, res.Error
	}
	return &link, nil
}

// UpdateLink обновляет данные ссылки в базе.
func (r *LinkRepository) UpdateLink(ctx context.Context, link *models.Link) (*models.Link, error) {
	res := r.Database.DB.WithContext(ctx).Clauses(clause.Returning{}).Updates(link)
	if res.Error != nil {
		log.Printf("[LinkRepository] Ошибка обновления ссылки (ID: %d): %v", link.ID, res.Error)
		return nil, fmt.Errorf("ошибка при обновлении ссылки в БД: %w", res.Error)
	}
	return link, nil
}

// DeleteLink удаляет ссылку по ID.
func (r *LinkRepository) DeleteLink(ctx context.Context, linkID uint) error {
	res := r.Database.DB.WithContext(ctx).Delete(&models.Link{}, linkID)
	if res.Error != nil {
		log.Printf("[LinkRepository] Ошибка удаления ссылки (ID: %d): %v", linkID, res.Error)
		return fmt.Errorf("ошибка при удалении ссылки из БД: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("ссылка с ID %d не найдена", linkID)
	}
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
		log.Printf("[LinkRepository] Ошибка подсчёта ссылок: %v", res.Error)
		return 0, fmt.Errorf("ошибка при подсчёте ссылок: %w", res.Error)
	}
	return count, nil
}

// FindLinkByID ищет ссылку по ID.
func (r *LinkRepository) FindLinkByID(ctx context.Context, linkID uint) (*models.Link, error) {
	var link models.Link
	res := r.Database.DB.WithContext(ctx).First(&link, linkID)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Printf("[LinkRepository] Ссылка с ID %d не найдена", linkID)
			return nil, nil
		}
		log.Printf("[LinkRepository] Ошибка при поиске ссылки (ID: %d): %v", linkID, res.Error)
		return nil, fmt.Errorf("ошибка при поиске ссылки: %w", res.Error)
	}
	return &link, nil
}
