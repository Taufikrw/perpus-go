package repository

import (
	"belajar-go/models"
	"context"

	"gorm.io/gorm"
)

type bookItemRepositoryImpl struct {
	db *gorm.DB
}

func NewBookItemRepository(db *gorm.DB) models.BookItemRepository {
	return &bookItemRepositoryImpl{db: db}
}

func (r *bookItemRepositoryImpl) FindByBookID(c context.Context, bookID string) ([]models.BookItem, error) {
	var bookItems []models.BookItem
	if err := r.db.Preload("Book.Category").Where("book_id = ?", bookID).Find(&bookItems).Error; err != nil {
		return nil, err
	}
	return bookItems, nil
}

func (r *bookItemRepositoryImpl) FindByID(c context.Context, id string) (*models.BookItem, error) {
	var bookItem models.BookItem
	if err := r.db.Preload("Book.Category").Where("id = ?", id).Take(&bookItem).Error; err != nil {
		return nil, err
	}
	return &bookItem, nil
}

func (r *bookItemRepositoryImpl) Create(c context.Context, bookItem *models.BookItem) error {
	db := GetDB(c, r.db)
	return db.Create(bookItem).Error
}

func (r *bookItemRepositoryImpl) Update(c context.Context, bookItem *models.BookItem) error {
	db := GetDB(c, r.db)
	return db.Updates(bookItem).Error
}

func (r *bookItemRepositoryImpl) Delete(c context.Context, bookItem *models.BookItem) error {
	db := GetDB(c, r.db)
	return db.Delete(bookItem).Error
}
