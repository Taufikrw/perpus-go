package resources

import "belajar-go/models"

type CategoryResource struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func FormatCategory(category models.BookCategory) CategoryResource {
	return CategoryResource{
		ID:        category.ID.String(),
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func FormatCategories(categories []models.BookCategory) []CategoryResource {
	var categoryResources []CategoryResource
	for _, category := range categories {
		categoryResources = append(categoryResources, FormatCategory(category))
	}
	return categoryResources
}
