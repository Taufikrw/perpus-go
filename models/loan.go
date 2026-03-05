package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Loan struct {
	BaseModel
	MemberID   uuid.UUID
	BookItemID uuid.UUID
	LoanDate   time.Time
	DueDate    time.Time
	ReturnDate *time.Time
	Status     string

	Member   Member   `gorm:"foreignKey:MemberID;references:ID"`
	BookItem BookItem `gorm:"foreignKey:BookItemID;references:ID"`
	Fine     *Fine    `gorm:"foreignKey:LoanID;references:ID"`
}

type LoanRepository interface {
	FindAll(c context.Context) ([]Loan, error)
	FindByID(c context.Context, id string) (*Loan, error)
	Create(c context.Context, loan *Loan) error
	Update(c context.Context, loan *Loan) error
	Delete(c context.Context, loan *Loan) error
}
