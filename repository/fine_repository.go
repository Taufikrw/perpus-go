package repository

import (
	"belajar-go/models"
	"context"

	"gorm.io/gorm"
)

type FineRepository interface {
	BaseRepository[models.Fine]
	FindByLoanID(c context.Context, loanID string) (*models.Fine, error)
}

type fineRepositoryImpl struct {
	BaseRepository[models.Fine]
	db *gorm.DB
}

func NewFineRepository(db *gorm.DB) FineRepository {
	return &fineRepositoryImpl{
		BaseRepository: NewBaseRepository[models.Fine](db),
		db:             db,
	}
}

func (r *fineRepositoryImpl) FindByLoanID(c context.Context, loanID string) (*models.Fine, error) {
	var fine models.Fine
	if err := r.db.WithContext(c).Preload("Loan.Member.User.Role").Preload("Loan.BookItem.Book.Category").Where("loan_id = ?", loanID).Take(&fine).Error; err != nil {
		return nil, err
	}
	return &fine, nil
}
