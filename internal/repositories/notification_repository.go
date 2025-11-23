package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/Dubjay18/scenee/internal/models"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *models.Notification) error
	GetByUserID(ctx context.Context, userID string) ([]models.Notification, error)
	MarkAsRead(ctx context.Context, id string) error
}

type GormNotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *GormNotificationRepository {
	return &GormNotificationRepository{db: db}
}

func (r *GormNotificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *GormNotificationRepository) GetByUserID(ctx context.Context, userID string) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}

func (r *GormNotificationRepository) MarkAsRead(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&models.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}
