package models

import (
	"context"
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

type FineRepository interface {
	FindAll(c context.Context) ([]Fine, error)
	FindByID(c context.Context, id string) (*Fine, error)
	Create(c context.Context, fine *Fine) error
	Update(c context.Context, fine *Fine) error
	Delete(c context.Context, fine *Fine) error
}
