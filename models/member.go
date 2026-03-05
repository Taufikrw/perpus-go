package models

import (
	"context"

	"github.com/google/uuid"
)

type Member struct {
	BaseModel
	MemberCode  string `gorm:"unique"`
	PhoneNumber string
	Address     string
	IsApproved  bool `gorm:"default:false"`
	UserID      uuid.UUID

	User  User   `gorm:"foreignKey:UserID;references:ID"`
	Loans []Loan `gorm:"foreignKey:MemberID;references:ID"`
}

type MemberRepository interface {
	FindAll(c context.Context) ([]Member, error)
	FindByID(c context.Context, id string) (*Member, error)
	FindByUserID(c context.Context, userID string) (*Member, error)
	Create(c context.Context, member *Member) error
	Update(c context.Context, member *Member) error
	Delete(c context.Context, member *Member) error
	IsMemberCodeExists(c context.Context, memberCode string, excludeID string) (bool, error)
}
