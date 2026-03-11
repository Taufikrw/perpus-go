package models

type BookCategory struct {
	BaseModel
	Name string

	Books []Book `gorm:"foreignKey:CategoryID;references:ID"`
}
