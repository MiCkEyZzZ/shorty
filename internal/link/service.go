package link

import (
	"context"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// Ошибки, связанные с работой с коращенными ссылками.
var (
	ErrLinkNotFound = errors.New("ссылка не найдена")
	ErrLinkCreation = errors.New("не удалось создать ссылку")
	ErrLinkUpdate   = errors.New("не удалось обновить ссылку")
	ErrLinkDeletion = errors.New("не удалось удалить ссылку")
	ErrLinkNotValid = errors.New("ссылка некорректна")
)

// LinkService - сервис для управления сокращёнными ссылками.
type LinkService struct {
	repo *LinkRepository
}

// NewLinkService создаёт новый экземпляр LinkService.
func NewLinkService(repo *LinkRepository) *LinkService {
	return &LinkService{repo: repo}
}

// Create создаёт новую сокращённую ссылку.
func (s *LinkService) Create(ctx context.Context, link *Link) (*Link, error) {
	newLink, err := s.repo.CreateLink(ctx, link)
	if err != nil {
		log.Printf("[LinkService] Ошибка при создании ссылки: %v", err)
		return nil, ErrLinkCreation
	}
	return newLink, nil
}

// Count получает число страниц.
func (s *LinkService) Count(ctx context.Context) int64 {
	return s.repo.Count(ctx)
}

// GetLinks получает список ссылок.
func (s *LinkService) GetLinks(ctx context.Context, limit, offset int) []Link {
	links := s.repo.GetLinks(ctx, limit, offset)
	return links
}

// GetByHash ищет сокращённую ссылку по её хешу.
func (s *LinkService) GetByHash(ctx context.Context, hash string) (*Link, error) {
	link, err := s.repo.GetLinkByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLinkNotFound
		}
		log.Printf("[LinkService] Ошибка получения ссылки по хешу %s: %v", hash, err)
		return nil, fmt.Errorf("не удалось найти ссылку с хешем %s: %w", hash, err)
	}
	return link, nil
}

// Update обновляет существующую сокращённую ссылку.
func (s *LinkService) Update(ctx context.Context, link *Link) (*Link, error) {
	updatedLink, err := s.repo.UpdateLink(ctx, link)
	if err != nil {
		log.Printf("[LinkService] Ошибка при обновлении ссылки (ID: %d): %v", link.ID, err)
		return nil, ErrLinkUpdate
	}
	return updatedLink, nil
}

// Delete удаляет сокращённую ссылку по её ID.
func (s *LinkService) Delete(ctx context.Context, id uint) error {
	err := s.repo.DeleteLink(ctx, id)
	if err != nil {
		log.Printf("[LinkService] Ошибка удаления ссылки (ID: %d): %v", id, err)
		return ErrLinkDeletion
	}
	return nil
}

// FindByID находит ссылку по её ID.
func (s *LinkService) FindByID(ctx context.Context, id uint) (*Link, error) {
	link, err := s.repo.FindLinkByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLinkNotFound
		}
		log.Printf("[LinkService] Ошибка при поиске ссылки (ID: %d): %v", id, err)
		return nil, fmt.Errorf("не удалось найти ссылку: %w", err)
	}
	return link, nil
}
