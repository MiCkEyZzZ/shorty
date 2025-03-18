package repository

import (
	"context"
	"log"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"shorty/internal/consts"
	"shorty/internal/models"
	"shorty/internal/payload"
	"shorty/pkg/db"
)

// StatRepository отвечает за операции с базой данных для сущности Stat.
type StatRepository struct {
	Database *db.DB
}

// NewStatRepository создает новый экземпляр StatRepository
func NewStatRepository(db *db.DB) *StatRepository {
	return &StatRepository{Database: db}
}

// AddClick инкрементирует количество кликов для ссылки на текущую дату
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
				log.Printf("[StatRepository] Ошибка при создании статистики: %v", err)
				return err
			}
		} else {
			// Ошибка запроса
			tx.Rollback()
			log.Printf("[StatRepository] Ошибка при поиске статистики: %v", err)
			return err
		}
	} else {
		// Если запись найдена, увеличиваем количество кликов
		if err := tx.Model(&models.Stat{}).
			Where("link_id = ? AND date = ?", linkID, currentDate).
			Update("clicks", gorm.Expr("clicks + ?", 1)).Error; err != nil {
			tx.Rollback()
			log.Printf("[StatRepository] Ошибка при обновлении статистики: %v", err)
			return err
		}
	}
	// Проверяем, не была ли уже откатана транзакция
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	return tx.Commit().Error
}

// GetStats получает статистику по дням или месяцам за заданный период
func (r *StatRepository) GetStats(ctx context.Context, by string, from, to time.Time) []payload.GetStatsResponse {
	var stats []payload.GetStatsResponse
	var selectQuery string

	// Определяем группировку.
	switch by {
	case consts.GroupByDay:
		selectQuery = "to_char(date, 'YYYY-MM-DD') as period, sum(clicks)"
	case consts.GroupByMonth:
		selectQuery = "to_char(date, 'YYYY-MM') as period, sum(clicks)"
	default:
		log.Printf("[StatRepository] Неверное значение для группировки: %s", by)
		return nil
	}
	r.Database.DB.
		Model(&models.Stat{}).
		WithContext(ctx).
		Select(selectQuery).
		Where("date BETWEEN ? AND ?", from, to).
		Group("period").
		Order("period").
		Scan(&stats)

	return stats
}
