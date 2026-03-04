package seeders

import (
	"belajar-go/models"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) {
	roles := []string{"admin", "librarian", "member"}

	for _, roleName := range roles {
		var role models.Role
		db.FirstOrCreate(&role, models.Role{Name: roleName})
	}

	fmt.Println("Seeding Roles selesai...")

	var adminRole models.Role
	db.Where("name = ?", "admin").First(&adminRole)

	var adminUser models.User
	if err := db.Where("email = ?", "admin@example.com").First(&adminUser).Error; err != nil {
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

		db.Create(&adminUser)
	} else {
		fmt.Println("Super Admin sudah ada di database, lewati seeding admin.")
	}
}
