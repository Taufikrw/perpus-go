package repository

import (
	"belajar-go/models"

	"gorm.io/gorm"
)

type BookRepository interface {
	BaseRepository[models.Book]
}

type bookRepositoryImpl struct {
	BaseRepository[models.Book]
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepositoryImpl{
		BaseRepository: NewBaseRepository[models.Book](db),
		db:             db,
	}
}
