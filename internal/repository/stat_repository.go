package repository

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"shorty/internal/common"
	"shorty/internal/models"
	"shorty/internal/payload"
	"shorty/pkg/db"
	"shorty/pkg/logger"
)

// StatRepository отвечает за операции с базой данных для сущности Stat.
type StatRepository struct {
	Database *db.DB
}

// NewStatRepository создаёт новый экземпляр StatRepository.
func NewStatRepository(db *db.DB) *StatRepository {
	return &StatRepository{Database: db}
}

// AddClick метод для инкрементации количества кликов для ссылки на текущую дату.
func (r *StatRepository) AddClick(ctx context.Context, linkID uint) error {
	var stat models.Stat
	currentDate := datatypes.Date(time.Now())

	// Используем транзакцию для атомарности обновления или вставки
	tx := r.Database.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Проверяем, есть ли уже запись за текущий день
	err := tx.Where("link_id = ? AND date = ?", linkID, currentDate).First(&stat).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Если записи нет — создаём новую
			stat = models.Stat{
				LinkID: linkID,
				Clicks: 1,
				Date:   currentDate,
			}
			if err := tx.Create(&stat).Error; err != nil {
				tx.Rollback()
				logger.Error("Ошибка при создании статистики", zap.Error(err))
				return err
			}
		} else {
			// Ошибка запроса
			tx.Rollback()
			logger.Error("Ошибка при поиске статистики", zap.Error(err))
			return err
		}
	} else {
		// Если запись найдена, увеличиваем количество кликов
		if err := tx.Model(&models.Stat{}).
			Where("link_id = ? AND date = ?", linkID, currentDate).
			Update("clicks", gorm.Expr("clicks + ?", 1)).Error; err != nil {
			tx.Rollback()
			logger.Error("Ошибка при обновлении статистики", zap.Error(err))
			return err
		}
		logger.Info("Количество кликов увеличено", zap.Uint("linkID", linkID), zap.Int("clicks", 1))
	}
	// Проверяем, не была ли уже откатана транзакция
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	return tx.Commit().Error
}

// GetStats метод для получения статистики по дням или месяцам за заданный период.
func (r *StatRepository) GetStats(ctx context.Context, by string, from, to time.Time) []payload.GetStatsResponse {
	var stats []payload.GetStatsResponse
	var selectQuery string

	// Определяем группировку.
	switch by {
	case common.GroupByDay:
		selectQuery = "to_char(date, 'YYYY-MM-DD') as period, sum(clicks)"
	case common.GroupByMonth:
		selectQuery = "to_char(date, 'YYYY-MM') as period, sum(clicks)"
	default:
		logger.Warn("Неверное значение для группировки", zap.String("groupBy", by))
		return nil
	}
	result := r.Database.DB.
		WithContext(ctx).
		Model(&models.Stat{}).
		Select(selectQuery).
		Where("date BETWEEN ? AND ?", from, to).
		Group("period").
		Order("period").
		Scan(&stats)
	if result.Error != nil {
		logger.Error("Ошибка при получении статистики", zap.String("groupBy", by), zap.Error(result.Error))
		return nil
	}
	logger.Info("Статистика успешно получена", zap.String("groupBy", by), zap.Int("statsCount", len(stats)))
	return stats
}

func (r *StatRepository) GetAllLinksStats(ctx context.Context, from, to time.Time) []payload.LinkStatsResponse {
	var stats []payload.LinkStatsResponse

	r.Database.DB.
		Model(&models.Stat{}).
		WithContext(ctx).
		Select(`
			links.id AS link_id,
			links.url AS url,
			COUNT(stats.id) AS total_clicks,
			MAX(stats.date) AS last_click_date,
			SUM(CASE WHEN links.is_blocked = true THEN 1 ELSE 0 END) AS blocked_count
		`).
		Joins("LEFT JOIN links ON stats.link_id = links.id").
		Where("stats.date BETWEEN ? AND ?", from, to).
		Group("links.id, links.url").
		Order("total_clicks DESC").
		Scan(&stats)

	logger.Info("Статистика по всем ссылкам получена", zap.Int("count", len(stats)))
	return stats
}
