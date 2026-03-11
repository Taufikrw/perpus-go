package repository

import (
	"belajar-go/utils"
	"context"

	"gorm.io/gorm"
)

type BaseRepository[T any] interface {
	GetAll(c context.Context, preloads ...string) ([]T, error)
	GetByID(c context.Context, id string, preloads ...string) (*T, error)
	GetAllTrashedData(c context.Context, preloads ...string) ([]T, error)
	GetTrashedDataByID(c context.Context, id string, preloads ...string) (*T, error)
	Create(c context.Context, entity *T) error
	Update(c context.Context, entity *T) error
	Delete(c context.Context, entity *T) error
	Restore(c context.Context, id string) error
}

type baseRepositoryImpl[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return &baseRepositoryImpl[T]{db: db}
}

func (r *baseRepositoryImpl[T]) GetAll(ctx context.Context, preloads ...string) ([]T, error) {
	var entities []T
	query := r.db.WithContext(ctx)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	err := query.Find(&entities).Error
	return entities, err
}

func (r *baseRepositoryImpl[T]) GetByID(ctx context.Context, id string, preloads ...string) (*T, error) {
	var entity T
	query := r.db.WithContext(ctx)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	err := query.Where("id = ?", id).Take(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *baseRepositoryImpl[T]) GetAllTrashedData(ctx context.Context, preloads ...string) ([]T, error) {
	var entities []T
	query := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL")
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	err := query.Find(&entities).Error
	return entities, err
}

func (r *baseRepositoryImpl[T]) GetTrashedDataByID(ctx context.Context, id string, preloads ...string) (*T, error) {
	var entity T
	query := r.db.WithContext(ctx).Unscoped().Where("id = ? AND deleted_at IS NOT NULL", id)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	err := query.Take(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *baseRepositoryImpl[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *baseRepositoryImpl[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *baseRepositoryImpl[T]) Delete(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Delete(entity).Error
}

func (r *baseRepositoryImpl[T]) Restore(ctx context.Context, id string) error {
	db := GetDB(ctx, r.db)
	result := db.Unscoped().Model(new(T)).Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return utils.NewNotFoundError("Entity not found")
	}
	return nil
}
