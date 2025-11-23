package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Dubjay18/scenee/internal/models"
)

type UserRepository interface {
	Upsert(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
}

type GormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Upsert(ctx context.Context, user *models.User) error {
	if user.ID.String() == "" && user.Email == "" && user.Username == "" {
		return errors.New("missing identifiers")
	}
	return r.db.WithContext(ctx).Where(models.User{Email: user.Email}).Assign(user).FirstOrCreate(user).Error
}

func (r *GormUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *GormUserRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
}
