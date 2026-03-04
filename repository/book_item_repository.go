package repository

import (
	"belajar-go/models"
	"context"

	"gorm.io/gorm"
)

type bookItemRepository struct {
	db *gorm.DB
}

func NewBookItemRepository(db *gorm.DB) models.BookItemRepositoryInterface {
	return &bookItemRepository{db: db}
}

func (r *bookItemRepository) FindByBookID(c context.Context, bookID string) ([]models.BookItem, error) {
	var bookItems []models.BookItem
	if err := r.db.Preload("Book.Category").Where("book_id = ?", bookID).Find(&bookItems).Error; err != nil {
		return nil, err
	}
	return bookItems, nil
}

func (r *bookItemRepository) FindByID(c context.Context, id string) (*models.BookItem, error) {
	var bookItem models.BookItem
	if err := r.db.Preload("Book.Category").Where("id = ?", id).Take(&bookItem).Error; err != nil {
		return nil, err
	}
	return &bookItem, nil
}

func (r *bookItemRepository) Create(c context.Context, bookItem *models.BookItem) error {
	db := GetDB(c, r.db)
	return db.Create(bookItem).Error
}

func (r *bookItemRepository) Update(c context.Context, bookItem *models.BookItem) error {
	db := GetDB(c, r.db)
	return db.Updates(bookItem).Error
}

func (r *bookItemRepository) Delete(c context.Context, bookItem *models.BookItem) error {
	db := GetDB(c, r.db)
	return db.Delete(bookItem).Error
}
