package models

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	BaseModel
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
	RoleID   uuid.UUID

	Role   Role    `gorm:"foreignKey:RoleID;references:ID"`
	Member *Member `gorm:"foreignKey:UserID;references:ID"`
}

type AuthRepositoryInterface interface {
	GetUserByEmail(c context.Context, email string) (*User, error)
	GetRoleByName(c context.Context, name string) (*Role, error)
	RegisterMemberTransaction(c context.Context, user *User, member *Member) (*Member, error)
}

type UserRepository interface {
	FindAll(c context.Context) ([]User, error)
	FindByID(c context.Context, id string) (*User, error)
	GetRoleByName(c context.Context, name string) (*Role, error)
	Create(c context.Context, user *User) error
	Update(c context.Context, user *User) error
	Delete(c context.Context, user *User) error
	IsEmailExists(c context.Context, email string, excludeID string) (bool, error)
	IsUsernameExists(c context.Context, username string, excludeID string) (bool, error)
}
