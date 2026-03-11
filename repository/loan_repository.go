package repository

import (
	"belajar-go/models"

	"gorm.io/gorm"
)

type LoanRepository interface {
	BaseRepository[models.Loan]
}

type loanRepositoryImpl struct {
	BaseRepository[models.Loan]
	db *gorm.DB
}

func NewLoanRepository(db *gorm.DB) LoanRepository {
	return &loanRepositoryImpl{
		BaseRepository: NewBaseRepository[models.Loan](db),
		db:             db,
	}
}
