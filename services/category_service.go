package services

import (
	"belajar-go/dto"
	"belajar-go/models"
	"belajar-go/utils"
	"context"
)

type CategoryService struct {
	repository models.CategoryRepositoryInterface
}

func NewCategoryService(repository models.CategoryRepositoryInterface) *CategoryService {
	return &CategoryService{repository: repository}
}

func (s *CategoryService) GetAllCategories(c context.Context) ([]models.BookCategory, error) {
	return s.repository.FindAll(c)
}

func (s *CategoryService) GetCategoryByID(c context.Context, id string) (*models.BookCategory, error) {
	category, err := s.repository.FindByID(c, id)
	if category == nil {
		return nil, utils.NewNotFoundError("Category not found")
	} else if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) CreateCategory(c context.Context, input dto.CategoryDTO) (*models.BookCategory, error) {
	category := models.BookCategory{
		Name: input.Name,
	}

	err := s.repository.Create(c, &category)
	return &category, err
}

func (s *CategoryService) UpdateCategory(c context.Context, id string, input dto.CategoryDTO) (*models.BookCategory, error) {
	category, err := s.GetCategoryByID(c, id)
	if err != nil {
		return nil, err
	}

	category.Name = input.Name
	err = s.repository.Update(c, category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) DeleteCategory(c context.Context, id string) error {
	category, err := s.GetCategoryByID(c, id)
	if err != nil {
		return err
	}
	return s.repository.Delete(c, category)
}
