package dto

type CategoryDTO struct {
	Name string `json:"name" binding:"required"`
}
