package repository

import (
	"belajar-go/models"
	"context"

	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) models.AuthRepositoryInterface {
	return &authRepository{db: db}
}

func (r *authRepository) GetUserByEmail(c context.Context, email string) (*models.User, error) {
	var user models.User

	if err := r.db.WithContext(c).Where("email = ?", email).Take(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetRoleByName(c context.Context, name string) (*models.Role, error) {
	var role models.Role
	if err := r.db.WithContext(c).Where("name = ?", name).Take(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *authRepository) RegisterMemberTransaction(c context.Context, user *models.User, member *models.Member) (*models.Member, error) {
	err := r.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		member.UserID = user.ID
		if err := tx.Create(member).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return member, nil
}
