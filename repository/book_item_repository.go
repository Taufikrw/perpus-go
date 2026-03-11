package repository

import (
	"belajar-go/models"
	"context"

	"gorm.io/gorm"
)

type BookItemRepository interface {
	BaseRepository[models.BookItem]
	FindByBookID(c context.Context, bookID string) ([]models.BookItem, error)
	IsInventoryCodeExists(c context.Context, inventoryCode string, excludeID string) (bool, error)
}

type bookItemRepositoryImpl struct {
	BaseRepository[models.BookItem]
	db *gorm.DB
}

func NewBookItemRepository(db *gorm.DB) BookItemRepository {
	return &bookItemRepositoryImpl{
		BaseRepository: NewBaseRepository[models.BookItem](db),
		db:             db,
	}
}

func (r *bookItemRepositoryImpl) FindByBookID(c context.Context, bookID string) ([]models.BookItem, error) {
	var bookItems []models.BookItem
	if err := r.db.Preload("Book.Category").Where("book_id = ?", bookID).Find(&bookItems).Error; err != nil {
		return nil, err
	}
	return bookItems, nil
}

func (r *bookItemRepositoryImpl) IsInventoryCodeExists(c context.Context, inventoryCode string, excludeID string) (bool, error) {
	var count int64
	db := GetDB(c, r.db).Model(&models.BookItem{}).Where("inventory_code = ?", inventoryCode)
	if excludeID != "" {
		db = db.Where("id != ?", excludeID)
	}
	if err := db.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
