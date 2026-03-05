package models

import (
	"context"
)

type BookCategory struct {
	BaseModel
	Name string

	Books []Book `gorm:"foreignKey:CategoryID;references:ID"`
}

type CategoryRepository interface {
	FindAll(c context.Context) ([]BookCategory, error)
	FindByID(c context.Context, id string) (*BookCategory, error)
	Create(c context.Context, category *BookCategory) error
	Update(c context.Context, category *BookCategory) error
	Delete(c context.Context, category *BookCategory) error
}
