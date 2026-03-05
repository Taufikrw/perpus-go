package repository

import (
	"belajar-go/models"
	"context"

	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) models.UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) FindAll(c context.Context) ([]models.User, error) {
	var users []models.User
	if err := r.db.WithContext(c).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepositoryImpl) FindByID(c context.Context, id string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(c).Where("id = ?", id).Take(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
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

func (r *userRepositoryImpl) Create(c context.Context, user *models.User) error {
	db := GetDB(c, r.db)
	return db.Create(user).Error
}

func (r *userRepositoryImpl) Update(c context.Context, user *models.User) error {
	db := GetDB(c, r.db)
	return db.Model(&models.User{}).Where("id = ?", user.ID).Updates(user).Error
}

func (r *userRepositoryImpl) Delete(c context.Context, user *models.User) error {
	db := GetDB(c, r.db)
	return db.Delete(user).Error
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
