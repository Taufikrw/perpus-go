package repository

import (
	"belajar-go/models"
	"context"

	"gorm.io/gorm"
)

type MemberRepositoryImpl struct {
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) models.MemberRepository {
	return &MemberRepositoryImpl{db: db}
}

func (r *MemberRepositoryImpl) FindAll(c context.Context) ([]models.Member, error) {
	var members []models.Member
	if err := r.db.WithContext(c).Preload("User.Role").Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *MemberRepositoryImpl) FindByID(c context.Context, id string) (*models.Member, error) {
	var member models.Member
	if err := r.db.WithContext(c).Preload("User.Role").Where("id = ?", id).Take(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *MemberRepositoryImpl) FindByUserID(c context.Context, userID string) (*models.Member, error) {
	var member models.Member
	if err := r.db.WithContext(c).Preload("User.Role").Where("user_id = ?", userID).Take(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *MemberRepositoryImpl) Create(c context.Context, member *models.Member) error {
	db := GetDB(c, r.db)
	return db.Create(member).Error
}

func (r *MemberRepositoryImpl) Update(ctx context.Context, member *models.Member) error {
	db := GetDB(ctx, r.db)
	return db.Model(&models.Member{}).Where("id = ?", member.ID).Updates(member).Error
}

func (r *MemberRepositoryImpl) Delete(ctx context.Context, member *models.Member) error {
	db := GetDB(ctx, r.db)
	return db.Delete(member).Error
}

func (r *MemberRepositoryImpl) IsMemberCodeExists(c context.Context, memberCode string, excludeID string) (bool, error) {
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
