package validators

import (
	"belajar-go/models"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type AppValidator struct {
	db *gorm.DB
}

func NewAppValidator(db *gorm.DB) *AppValidator {
	return &AppValidator{db: db}
}

func (v *AppValidator) ValidationUniqueEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	var count int64

	v.db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	return count == 0
}

func (v *AppValidator) ValidationUniqueUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	var count int64

	v.db.Model(&models.User{}).Where("username = ?", username).Count(&count)
	return count == 0
}

func (v *AppValidator) ValidationUniqueMemberCode(fl validator.FieldLevel) bool {
	memberCode := fl.Field().String()
	var count int64

	v.db.Model(&models.Member{}).Where("member_code = ?", memberCode).Count(&count)

	return count == 0
}

func (v *AppValidator) ValidationUniqueInventoryCode(fl validator.FieldLevel) bool {
	inventoryCode := fl.Field().String()
	var count int64

	v.db.Model(&models.BookItem{}).Where("inventory_code = ?", inventoryCode).Count(&count)

	return count == 0
}
