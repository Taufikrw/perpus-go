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

	db := config.ConnectDatabase()

	appValidator := validators.NewAppValidator(db)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("unique_email", appValidator.ValidationUniqueEmail)
		v.RegisterValidation("unique_username", appValidator.ValidationUniqueUsername)
		v.RegisterValidation("unique_member_code", appValidator.ValidationUniqueMemberCode)
		v.RegisterValidation("unique_inventory_code", appValidator.ValidationUniqueInventoryCode)
	}

	seeders.SeedDatabase(db)
	appCtrl := routes.InitDependency(db)
	r := routes.SetupRouter(appCtrl)

	r.Run(":8000")
}
