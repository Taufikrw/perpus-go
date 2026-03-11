package repository

import (
	"belajar-go/models"
	"context"

	"gorm.io/gorm"
)

type UserRepository interface {
	BaseRepository[models.User]
	GetUserByEmail(c context.Context, email string) (*models.User, error)
	GetRoleByName(c context.Context, name string) (*models.Role, error)
	IsEmailExists(c context.Context, email string, excludeID string) (bool, error)
	IsUsernameExists(c context.Context, username string, excludeID string) (bool, error)
}

type userRepositoryImpl struct {
	BaseRepository[models.User]
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{
		BaseRepository: NewBaseRepository[models.User](db),
		db:             db,
	}
}

func (r *userRepositoryImpl) GetUserByEmail(c context.Context, email string) (*models.User, error) {
	var user models.User

	if err := r.db.WithContext(c).Where("email = ?", email).Take(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) GetRoleByName(c context.Context, name string) (*models.Role, error) {
	var role models.Role
	if err := r.db.WithContext(c).Where("name = ?", name).Take(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *userRepositoryImpl) IsEmailExists(c context.Context, email string, excludeID string) (bool, error) {
	var count int64
	db := GetDB(c, r.db).Model(&models.User{}).Where("email = ?", email)
	if excludeID != "" {
		db = db.Where("id != ?", excludeID)
	}
	if err := db.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepositoryImpl) IsUsernameExists(c context.Context, username string, excludeID string) (bool, error) {
	var count int64
	db := GetDB(c, r.db).Model(&models.User{}).Where("username = ?", username)
	if excludeID != "" {
		db = db.Where("id != ?", excludeID)
	}
	if err := db.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
