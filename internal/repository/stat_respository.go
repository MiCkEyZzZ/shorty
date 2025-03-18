package repository

import (
	"context"
	"log"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"shorty/internal/models"
	"shorty/pkg/db"
)

type StatRepository struct {
	Database *db.DB
}

func NewStatRepository(db *db.DB) *StatRepository {
	return &StatRepository{Database: db}
}

func (r *StatRepository) AddClick(ctx context.Context, linkID uint) error {
	var stat models.Stat
	currentDate := datatypes.Date(time.Now())

	// Используем транзакцию для атомарности.
	tx := r.Database.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.First(&stat, "link_id = ? AND date = ?", linkID, currentDate).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Если записи нет, то создаём.
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
			tx.Rollback()
			log.Printf("[StatRepository] Ошибка при поиске статистики: %v", err)
			return err
		}
	} else {
		// Если запись найдена, то обновляем.
		stat.Clicks += 1
		if err := tx.Save(&stat).Error; err != nil {
			tx.Rollback()
			log.Printf("[StatRepository] Ошибка при обновлении статистики: %v", err)
			return err
		}
	}
	return tx.Commit().Error
}
