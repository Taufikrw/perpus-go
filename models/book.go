package models

import (
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
