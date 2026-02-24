package models

import "github.com/google/uuid"

type User struct {
	BaseModel
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
	RoleID   uuid.UUID

	Role   Role    `gorm:"foreignKey:RoleID;references:ID"`
	Member *Member `gorm:"foreignKey:UserID;references:ID"`
}
