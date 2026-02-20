package main

import (
	"belajar-go/config"
	"belajar-go/routes"
	"belajar-go/seeders"
	"belajar-go/validators"
	"log"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan environment OS bawaan")
	}

	config.ConnectDatabase()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("unique_email", validators.ValidationUniqueEmail)
		v.RegisterValidation("unique_username", validators.ValidationUniqueUsername)
		v.RegisterValidation("unique_member_code", validators.ValidationUniqueMemberCode)
	}

	seeders.SeedRoles()
	r := routes.SetupRouter()

	r.Run(":8000")
}
