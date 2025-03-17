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

// LinkRepository предоставляет методы для работы с сокращёнными ссылками в БД.
type LinkRepository struct {
	Database *db.DB
}

// NewLinkRepository создаёт новый экземпляр LinkRepository.
func NewLinkRepository(db *db.DB) *LinkRepository {
	return &LinkRepository{
		Database: db,
	}
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

// GetLinkByHash ищет ссылку в базе данных по её хешу.
func (r *LinkRepository) GetLinkByHash(ctx context.Context, hash string) (*models.Link, error) {
	var link models.Link
	res := r.Database.DB.WithContext(ctx).First(&link, "hash = ?", hash)
	if res.Error != nil {
		log.Printf("[LinkRepository] Ошибка поиска ссылки: %v", res.Error)
		return nil, res.Error
	}

	return &link, nil
}

// UpdateLink обновляет данные ссылки в базе данных.
func (r *LinkRepository) UpdateLink(ctx context.Context, link *models.Link) (*models.Link, error) {
	res := r.Database.DB.WithContext(ctx).Clauses(clause.Returning{}).Updates(link)
	if res.Error != nil {
		log.Printf("[LinkRepository] Ошибка обновления ссылки (ID: %d): %v", link.ID, res.Error)
		return nil, fmt.Errorf("ошибка при обновлении ссылки в БД: %w", res.Error)
	}
	return link, nil
}

// DeleteLink удаляет ссылку из базы данных по её ID.
func (r *LinkRepository) DeleteLink(ctx context.Context, id uint) error {
	res := r.Database.DB.WithContext(ctx).Delete(&models.Link{}, id)
	if res.Error != nil {
		log.Printf("[LinkRepository] Ошибка удаления ссылки (ID: %d): %v", id, res.Error)
		return fmt.Errorf("ошибка при удалении ссылки из БД: %w", res.Error)
	}
	return nil
}

// FindLinkByID находит ссылку по её ID.
func (r *LinkRepository) FindLinkByID(ctx context.Context, id uint) (*models.Link, error) {
	var link models.Link
	res := r.Database.DB.WithContext(ctx).First(&link, id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Printf("[LinkRepository] Ссылка с ID %d не найдена", id)
			return nil, fmt.Errorf("ссылка с ID %d не найдена", id)
		}
		log.Printf("[LinkRepository] Ошибка при поиске ссылки (ID: %d): %v", id, res.Error)
		return nil, fmt.Errorf("ошибка при поиске ссылки: %w", res.Error)
	}

	return &link, nil
}
