package repository

import (
	"belajar-go/models"
	"context"

	"gorm.io/gorm"
)

type bookRepositoryImpl struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) models.BookRepository {
	return &bookRepositoryImpl{db: db}
}

func (r *bookRepositoryImpl) FindAll(c context.Context) ([]models.Book, error) {
	var books []models.Book
	if err := r.db.WithContext(c).Preload("Category").Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *bookRepositoryImpl) FindByID(c context.Context, id string) (*models.Book, error) {
	var book models.Book
	if err := r.db.WithContext(c).Preload("Category").Where("id = ?", id).Take(&book).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepositoryImpl) Create(c context.Context, book *models.Book) error {
	return r.db.WithContext(c).Create(book).Error
}

func (r *bookRepositoryImpl) Update(c context.Context, book *models.Book) error {
	return r.db.WithContext(c).Updates(book).Error
}

func (r *bookRepositoryImpl) Delete(c context.Context, book *models.Book) error {
	return r.db.WithContext(c).Delete(book).Error
}
