package dto

type CreateBookItemDTO struct {
	BookID        string `json:"book_id" binding:"required,uuid4"`
	InventoryCode string `json:"inventory_code" binding:"required,unique_inventory_code"`
	Condition     string `json:"condition" binding:"required,oneof=new good damaged"`
}

type UpdateBookItemDTO struct {
	BookID        string `json:"book_id" binding:"required,uuid4"`
	InventoryCode string `json:"inventory_code" binding:"required"`
	Condition     string `json:"condition" binding:"required,oneof=new good damaged"`
}
