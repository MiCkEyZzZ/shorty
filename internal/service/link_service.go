package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"

	"shorty/internal/models"
	"shorty/internal/repository"
)

var (
	ErrLinkNotFound = errors.New("ссылка не найдена")
	ErrLinkCreation = errors.New("не удалось создать ссылку")
	ErrLinkUpdate   = errors.New("не удалось обновить ссылку")
	ErrLinkDeletion = errors.New("не удалось удалить ссылку")
	ErrLinkNotValid = errors.New("ссылка некорректна")
)

type LinkService struct {
	Repo *repository.LinkRepository
}

// NewLinkService создаёт новый экземпляр LinkService
func NewLinkService(repo *repository.LinkRepository) *LinkService {
	return &LinkService{Repo: repo}
}

// Create создаёт новую ссылку
func (s *LinkService) Create(ctx context.Context, link *models.Link) (*models.Link, error) {
	newLink, err := s.Repo.CreateLink(ctx, link)
	if err != nil {
		log.Printf("[LinkService] Ошибка при создании ссылки: %v", err)
		return nil, ErrLinkCreation
	}
	return newLink, nil
}

// GetAll возвращает список ссылок с пагинацией
func (s *LinkService) GetAll(ctx context.Context, limit, offset int) ([]models.Link, error) {
	links, err := s.Repo.GetLinks(ctx, limit, offset)
	if err != nil {
		log.Printf("[LinkService] Ошибка при получении списка ссылок: %v", err)
		return nil, err
	}
	return links, nil
}

// GetByHash ищет ссылку по хешу
func (s *LinkService) GetByHash(ctx context.Context, hash string) (*models.Link, error) {
	link, err := s.Repo.GetLinkHash(ctx, hash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLinkNotFound
		}
		log.Printf("[LinkService] Ошибка получения ссылки по хешу %s: %v", hash, err)
		return nil, fmt.Errorf("не удалось найти ссылку с хешем %s: %w", hash, err)
	}
	return link, nil
}

// Update обновляет ссылку
func (s *LinkService) Update(ctx context.Context, link *models.Link) (*models.Link, error) {
	updatedLink, err := s.Repo.UpdateLink(ctx, link)
	if err != nil {
		log.Printf("[LinkService] Ошибка при обновлении ссылки (ID: %d): %v", link.ID, err)
		return nil, ErrLinkUpdate
	}
	return updatedLink, nil
}

// Delete удаляет ссылку по ID
func (s *LinkService) Delete(ctx context.Context, linkID uint) error {
	err := s.Repo.DeleteLink(ctx, linkID)
	if err != nil {
		log.Printf("[LinkService] Ошибка удаления ссылки (ID: %d): %v", linkID, err)
		return ErrLinkDeletion
	}
	return nil
}

// Count возвращает количество ссылок в базе
func (s *LinkService) Count(ctx context.Context) (int64, error) {
	res, err := s.Repo.CountLink(ctx)
	if err != nil {
		log.Printf("[LinkService] Ошибка при подсчёте ссылок: %v", err)
		return 0, err
	}
	return res, nil
}

// FindByID ищет ссылку по ID
func (s *LinkService) FindByID(ctx context.Context, linkID uint) (*models.Link, error) {
	link, err := s.Repo.FindLinkByID(ctx, linkID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLinkNotFound
		}
		log.Printf("[LinkService] Ошибка при поиске ссылки (ID: %d): %v", linkID, err)
		return nil, fmt.Errorf("не удалось найти ссылку: %w", err)
	}
	return link, nil
}
