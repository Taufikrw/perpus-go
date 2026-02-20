package models

import "github.com/google/uuid"

type BookItem struct {
	BaseModel
	BookID        uuid.UUID
	InventoryCode string `gorm:"unique"`
	Condition     string
	Status        string

	Book  Book   `gorm:"foreignKey:BookID;references:ID"`
	Loans []Loan `gorm:"foreignKey:BookItemID;references:ID"`
}
