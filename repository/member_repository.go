package repository

import (
	"belajar-go/models"
	"context"

	"gorm.io/gorm"
)

type MemberRepository interface {
	BaseRepository[models.Member]
	FindByUserID(c context.Context, userID string) (*models.Member, error)
	IsMemberCodeExists(c context.Context, memberCode string, excludeID string) (bool, error)
}

type memberRepositoryImpl struct {
	BaseRepository[models.Member]
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) MemberRepository {
	return &memberRepositoryImpl{
		BaseRepository: NewBaseRepository[models.Member](db),
		db:             db,
	}
}

func (r *memberRepositoryImpl) FindByUserID(c context.Context, userID string) (*models.Member, error) {
	var member models.Member
	if err := r.db.WithContext(c).Preload("User.Role").Where("user_id = ?", userID).Take(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *memberRepositoryImpl) IsMemberCodeExists(c context.Context, memberCode string, excludeID string) (bool, error) {
	var count int64
	db := GetDB(c, r.db).Model(&models.Member{}).Where("member_code = ?", memberCode)
	if excludeID != "" {
		db = db.Where("id != ?", excludeID)
	}
	if err := db.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
