package repository

import (
	"belajar-go/models"
	"context"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) models.UserRepositoryInterface {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll(c context.Context) ([]models.User, error) {
	var users []models.User
	if err := r.db.WithContext(c).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) FindByID(c context.Context, id string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(c).Where("id = ?", id).Take(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetRoleByName(c context.Context, name string) (*models.Role, error) {
	var role models.Role
	if err := r.db.WithContext(c).Where("name = ?", name).Take(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *userRepository) Create(c context.Context, user *models.User) error {
	db := GetDB(c, r.db)
	return db.Create(user).Error
}

func (r *userRepository) Update(c context.Context, user *models.User) error {
	db := GetDB(c, r.db)
	return db.Model(&models.User{}).Where("id = ?", user.ID).Updates(user).Error
}

func (r *userRepository) Delete(c context.Context, user *models.User) error {
	db := GetDB(c, r.db)
	return db.Delete(user).Error
}

func (r *userRepository) IsEmailExists(c context.Context, email string, excludeID string) (bool, error) {
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

func (r *userRepository) IsUsernameExists(c context.Context, username string, excludeID string) (bool, error) {
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
