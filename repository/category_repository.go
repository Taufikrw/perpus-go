package repository

import (
	"belajar-go/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	BaseRepository[models.BookCategory]
}

type categoryRepositoryImpl struct {
	BaseRepository[models.BookCategory]
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepositoryImpl{
		BaseRepository: NewBaseRepository[models.BookCategory](db),
		db:             db,
	}
}
