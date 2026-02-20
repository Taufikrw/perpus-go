package models

import (
	"time"

	"github.com/google/uuid"
)

type Loan struct {
	BaseModel
	UserID     uuid.UUID
	BookItemID uuid.UUID
	LoanDate   time.Time
	DueDate    time.Time
	ReturnDate *time.Time
	Status     string

	User     User     `gorm:"foreignKey:UserID;references:ID"`
	BookItem BookItem `gorm:"foreignKey:BookItemID;references:ID"`
	Fine     *Fine    `gorm:"foreignKey:LoanID;references:ID"`
}
