package repository

import (
	"belajar-go/models"
	"context"

	"gorm.io/gorm"
)

type fineRepositoryImpl struct {
	db *gorm.DB
}

func NewFineRepository(db *gorm.DB) models.FineRepository {
	return &fineRepositoryImpl{db: db}
}

func (r *fineRepositoryImpl) FindAll(c context.Context) ([]models.Fine, error) {
	var fines []models.Fine
	if err := r.db.WithContext(c).Preload("Loan.Member.User.Role").Preload("Loan.BookItem.Book.Category").Find(&fines).Error; err != nil {
		return nil, err
	}
	return fines, nil
}

func (r *fineRepositoryImpl) FindByID(c context.Context, id string) (*models.Fine, error) {
	var fine models.Fine
	if err := r.db.WithContext(c).Preload("Loan.Member.User.Role").Preload("Loan.BookItem.Book.Category").Where("id = ?", id).Take(&fine).Error; err != nil {
		return nil, err
	}
	return &fine, nil
}

func (r *fineRepositoryImpl) Create(c context.Context, fine *models.Fine) error {
	db := GetDB(c, r.db)
	return db.Create(fine).Error
}

func (r *fineRepositoryImpl) Update(c context.Context, fine *models.Fine) error {
	db := GetDB(c, r.db)
	return db.Updates(fine).Error
}

func (r *fineRepositoryImpl) Delete(c context.Context, fine *models.Fine) error {
	db := GetDB(c, r.db)
	return db.Delete(fine).Error
}
