package repository

import (
	"belajar-go/models"
	"context"

	"gorm.io/gorm"
)

type loanRepositoryImpl struct {
	db *gorm.DB
}

func NewLoanRepository(db *gorm.DB) models.LoanRepository {
	return &loanRepositoryImpl{db: db}
}

func (r *loanRepositoryImpl) FindAll(c context.Context) ([]models.Loan, error) {
	var loans []models.Loan
	if err := r.db.WithContext(c).Preload("Member.User.Role").Preload("BookItem.Book.Category").Preload("Fine").Find(&loans).Error; err != nil {
		return nil, err
	}
	return loans, nil
}

func (r *loanRepositoryImpl) FindByID(c context.Context, id string) (*models.Loan, error) {
	var loan models.Loan
	if err := r.db.WithContext(c).Preload("Member.User.Role").Preload("BookItem.Book.Category").Preload("Fine").Where("id = ?", id).Take(&loan).Error; err != nil {
		return nil, err
	}
	return &loan, nil
}

func (r *loanRepositoryImpl) Create(c context.Context, loan *models.Loan) error {
	db := GetDB(c, r.db)
	return db.Create(loan).Error
}

func (r *loanRepositoryImpl) Update(c context.Context, loan *models.Loan) error {
	db := GetDB(c, r.db)
	return db.Updates(loan).Error
}

func (r *loanRepositoryImpl) Delete(c context.Context, loan *models.Loan) error {
	db := GetDB(c, r.db)
	return db.Delete(loan).Error
}
