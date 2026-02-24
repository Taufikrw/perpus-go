package seeders

import (
	"belajar-go/config"
	"belajar-go/models"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func SeedDatabase() {
	roles := []string{"admin", "librarian", "member"}

	for _, roleName := range roles {
		var role models.Role
		config.DB.FirstOrCreate(&role, models.Role{Name: roleName})
	}

	fmt.Println("Seeding Roles selesai...")

	var adminRole models.Role
	config.DB.Where("name = ?", "admin").First(&adminRole)

	var adminUser models.User
	if err := config.DB.Where("email = ?", "admin@example.com").First(&adminUser).Error; err != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("asdf1234"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Gagal hash password seeder")
		}

		adminUser = models.User{
			Username: "Super Admin",
			Email:    "admin@example.com",
			Password: string(hashedPassword),
			RoleID:   adminRole.ID,
		}

		config.DB.Create(&adminUser)
	} else {
		fmt.Println("Super Admin sudah ada di database, lewati seeding admin.")
	}
}
