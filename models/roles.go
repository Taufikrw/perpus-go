package models

type Role struct {
	BaseModel
	Name string

	Users []User `gorm:"foreignKey:RoleID;references:ID"`
}
