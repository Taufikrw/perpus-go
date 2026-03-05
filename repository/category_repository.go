package repository

import (
	"belajar-go/models"
	"context"

	"gorm.io/gorm"
)

type categoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) models.CategoryRepository {
	return &categoryRepositoryImpl{db: db}
}

func (r *categoryRepositoryImpl) FindAll(c context.Context) ([]models.BookCategory, error) {
	var categories []models.BookCategory
	if err := r.db.WithContext(c).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepositoryImpl) FindByID(c context.Context, id string) (*models.BookCategory, error) {
	var category models.BookCategory
	if err := r.db.WithContext(c).Where("id = ?", id).Take(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepositoryImpl) Create(c context.Context, category *models.BookCategory) error {
	return r.db.WithContext(c).Create(category).Error
}

func (r *categoryRepositoryImpl) Update(ctx context.Context, category *models.BookCategory) error {
	return r.db.WithContext(ctx).Updates(category).Error
}

func (r *categoryRepositoryImpl) Delete(ctx context.Context, category *models.BookCategory) error {
	return r.db.WithContext(ctx).Delete(category).Error
}
