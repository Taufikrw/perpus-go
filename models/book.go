package models

import (
	"context"

	"github.com/google/uuid"
)

type Book struct {
	BaseModel
	CategoryID uuid.UUID
	Title      string
	Author     string
	Publisher  string
	Isbn       string
	Year       int
	Synopsis   string

	Category  BookCategory `gorm:"foreignKey:CategoryID;references:ID"`
	BookItems []BookItem   `gorm:"foreignKey:BookID;references:ID"`
}

type BookRepositoryInterface interface {
	FindAll(c context.Context) ([]Book, error)
	FindByID(c context.Context, id string) (*Book, error)
	Create(c context.Context, book *Book) error
	Update(c context.Context, book *Book) error
	Delete(c context.Context, book *Book) error
}
