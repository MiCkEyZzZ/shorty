package service

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"shorty/internal/models"
	"shorty/internal/repository"
	"shorty/pkg/logger"
)

var (
	ErrLinkNotFound = errors.New("ссылка не найдена")
	ErrLinkCreation = errors.New("не удалось создать ссылку")
	ErrLinkUpdate   = errors.New("не удалось обновить ссылку")
	ErrLinkDeletion = errors.New("не удалось удалить ссылку")
	ErrLinkNotValid = errors.New("ссылка некорректна")
)

// LinkService предоставляет методы для работы с ссылками.
type LinkService struct {
	Repo repository.LinkRepo
}

// NewLinkService создаёт новый экземпляр LinkService
func NewLinkService(repo repository.LinkRepo) *LinkService {
	return &LinkService{Repo: repo}
}

// Create создаёт новую ссылку
func (s *LinkService) Create(ctx context.Context, link *models.Link) (*models.Link, error) {
	newLink, err := s.Repo.CreateLink(ctx, link)
	if err != nil {
		logger.Error("Ошибка при создании ссылки", zap.Error(err))
		return nil, ErrLinkCreation
	}
	logger.Info("Ссылка успешно создана", zap.Uint("id", newLink.ID), zap.String("hash", newLink.Hash))
	return newLink, nil
}

// GetAll возвращает список ссылок с пагинацией
func (s *LinkService) GetAll(ctx context.Context, limit, offset int) ([]models.Link, error) {
	links, err := s.Repo.GetLinks(ctx, limit, offset)
	if err != nil {
		logger.Error("Ошибка при получении списка ссылок", zap.Error(err))
		return nil, err
	}
	logger.Info("Получен список ссылок", zap.Int("limit", limit), zap.Int("offset", offset))
	return links, nil
}

// GetByHash ищет ссылку по хешу
func (s *LinkService) GetByHash(ctx context.Context, hash string) (*models.Link, error) {
	link, err := s.Repo.GetLinkHash(ctx, hash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("Ссылка не найдена", zap.String("hash", hash))
			return nil, ErrLinkNotFound
		}
		logger.Error("Ошибка получения ссылки по хешу", zap.String("hash", hash), zap.Error(err))
		return nil, fmt.Errorf("не удалось найти ссылку с хешем %s: %w", hash, err)
	}
	logger.Info("Ссылка найдена по хешу", zap.String("hash", hash), zap.Uint("id", link.ID))
	return link, nil
}

// Update обновляет ссылку
func (s *LinkService) Update(ctx context.Context, link *models.Link) (*models.Link, error) {
	updatedLink, err := s.Repo.UpdateLink(ctx, link)
	if err != nil {
		logger.Error("Ошибка при обновлении ссылки", zap.Uint("id", link.ID), zap.Error(err))
		return nil, ErrLinkUpdate
	}
	logger.Info("Ссылка успешно обновлена", zap.Uint("id", updatedLink.ID), zap.String("hash", updatedLink.Hash))
	return updatedLink, nil
}

// Delete удаляет ссылку по ID
func (s *LinkService) Delete(ctx context.Context, linkID uint) error {
	err := s.Repo.DeleteLink(ctx, linkID)
	if err != nil {
		logger.Error("Ошибка удаления ссылки", zap.Uint("id", linkID), zap.Error(err))
		return ErrLinkDeletion
	}
	logger.Info("Ссылка успешно удалена", zap.Uint("id", linkID))
	return nil
}

// Count возвращает количество ссылок в базе
func (s *LinkService) Count(ctx context.Context) (int64, error) {
	res, err := s.Repo.CountLinks(ctx)
	if err != nil {
		logger.Error("Ошибка при подсчёте ссылок", zap.Error(err))
		return 0, err
	}
	logger.Info("Подсчитано количество ссылок", zap.Int64("count", res))
	return res, nil
}

// FindByID ищет ссылку по ID
func (s *LinkService) FindByID(ctx context.Context, linkID uint) (*models.Link, error) {
	link, err := s.Repo.FindLinkByID(ctx, linkID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("Ссылка не найдена", zap.Uint("id", linkID))
			return nil, ErrLinkNotFound
		}
		logger.Error("Ошибка при поиске ссылки", zap.Uint("id", linkID), zap.Error(err))
		return nil, fmt.Errorf("не удалось найти ссылку: %w", err)
	}
	logger.Info("Ссылка найдена по ID", zap.Uint("id", linkID), zap.String("hash", link.Hash))
	return link, nil
}

// Block блокирует ссылку
func (s *LinkService) Block(ctx context.Context, linkID uint) (*models.Link, error) {
	link, err := s.Repo.FindLinkByID(ctx, linkID)
	if err != nil {
		return nil, err
	}
	if link == nil {
		return nil, ErrLinkNotFound
	}

	link.IsBlocked = true
	updatedLink, err := s.Repo.BlockLink(ctx, link)
	if err != nil {
		logger.Error("Ошибка при блокировки ссылки", zap.Uint("id", linkID), zap.Error(err))
		return nil, ErrLinkUpdate
	}

	return updatedLink, nil
}

// UnBlock разблокирует ссылку
func (s *LinkService) UnBlock(ctx context.Context, linkID uint) (*models.Link, error) {
	link, err := s.Repo.FindLinkByID(ctx, linkID)
	if err != nil {
		return nil, err
	}
	if link == nil {
		return nil, ErrLinkNotFound
	}

	link.IsBlocked = false
	updatedLink, err := s.Repo.UnBlockLink(ctx, link)
	if err != nil {
		logger.Error("Ошибка при снятии блокировки с ссылки", zap.Uint("id", linkID), zap.Error(err))
		return nil, ErrLinkUpdate
	}

	return updatedLink, nil
}

// GetBlockedLinksCount метод для получения количества заблокированных ссылок.
func (s *LinkService) GetBlockedLinksCount(ctx context.Context) (int64, error) {
	count, err := s.Repo.GetBlockedLinksCount(ctx)
	if err != nil {
		logger.Error("Ошибка при получении количества заблокированных ссылок", zap.Error(err))
		return 0, err
	}
	logger.Info("Количество заблокированных ссылок получено", zap.Int64("count", count))
	return count, nil
}

// GetDeletedLinksCount метод для получения количества удалённых ссылок.
func (s *LinkService) GetDeletedLinksCount(ctx context.Context) (int64, error) {
	count, err := s.Repo.GetBlockedLinksCount(ctx)
	if err != nil {
		logger.Error("Ошибка при получении количества удалённых ссылок", zap.Error(err))
		return 0, err
	}
	logger.Info("Количество удалённых ссылок получено", zap.Int64("count", count))
	return count, nil
}

// GetTotalLinks метод для получения общего количества ссылок.
func (s *LinkService) GetTotalLinks(ctx context.Context) (int64, error) {
	count, err := s.Repo.GetTotalLinks(ctx)
	if err != nil {
		logger.Error("Ошибка при получении количества созданных ссылок", zap.Error(err))
		return 0, err
	}
	logger.Info("Количество созданных ссылок получено", zap.Int64("count", count))
	return count, nil
}
