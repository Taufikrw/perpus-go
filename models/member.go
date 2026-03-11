package models

import (
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
