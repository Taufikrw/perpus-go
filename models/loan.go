package models

import (
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
