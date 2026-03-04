package repository

import (
	"belajar-go/models"
	"context"

	"gorm.io/gorm"
)

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) models.BookRepositoryInterface {
	return &bookRepository{db: db}
}

func (r *bookRepository) FindAll(c context.Context) ([]models.Book, error) {
	var books []models.Book
	if err := r.db.WithContext(c).Preload("Category").Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *bookRepository) FindByID(c context.Context, id string) (*models.Book, error) {
	var book models.Book
	if err := r.db.WithContext(c).Preload("Category").Where("id = ?", id).Take(&book).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) Create(c context.Context, book *models.Book) error {
	return r.db.WithContext(c).Create(book).Error
}

func (r *bookRepository) Update(c context.Context, book *models.Book) error {
	return r.db.WithContext(c).Updates(book).Error
}

func (r *bookRepository) Delete(c context.Context, book *models.Book) error {
	return r.db.WithContext(c).Delete(book).Error
}
