package models

import (
	"context"

	"github.com/google/uuid"
)

type BookItem struct {
	BaseModel
	BookID        uuid.UUID
	InventoryCode string `gorm:"unique"`
	Condition     string
	Status        string

	Book  Book   `gorm:"foreignKey:BookID;references:ID"`
	Loans []Loan `gorm:"foreignKey:BookItemID;references:ID"`
}

type BookItemRepositoryInterface interface {
	FindByBookID(c context.Context, bookID string) ([]BookItem, error)
	FindByID(c context.Context, id string) (*BookItem, error)
	Create(c context.Context, bookItem *BookItem) error
	Update(c context.Context, bookItem *BookItem) error
	Delete(c context.Context, bookItem *BookItem) error
}
