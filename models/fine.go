package models

import (
	"time"

	"github.com/google/uuid"
)

type Fine struct {
	BaseModel
	LoanID uuid.UUID
	Amount float64
	PaidAt *time.Time

	Loan Loan `gorm:"foreignKey:LoanID;references:ID"`
}
