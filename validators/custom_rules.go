package validators

import (
	"belajar-go/config"
	"belajar-go/models"

	"github.com/go-playground/validator/v10"
)

func ValidationUniqueEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	var count int64

	config.DB.Model(&models.User{}).Where("email = ?", email).Count(&count)

	return count == 0
}

func ValidationUniqueUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	var count int64

	config.DB.Model(&models.User{}).Where("username = ?", username).Count(&count)

	return count == 0
}

func ValidationUniqueMemberCode(fl validator.FieldLevel) bool {
	memberCode := fl.Field().String()
	var count int64

	config.DB.Model(&models.Member{}).Where("member_code = ?", memberCode).Count(&count)

	return count == 0
}
