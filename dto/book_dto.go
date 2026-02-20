package dto

type CreateBookDTO struct {
	Title      string `json:"title" binding:"required,min=5"`
	Author     string `json:"author" binding:"required"`
	Year       int    `json:"year" binding:"required,number,gt=1900"`
	Publisher  string `json:"publisher"`
	Isbn       string `json:"isbn"`
	Synopsis   string `json:"synopsis"`
	CategoryID string `json:"category_id" binding:"required,uuid4"`
}
