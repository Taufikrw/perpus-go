package seeders

import (
	"belajar-go/models"
	"log"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	roles := []string{"admin", "member", "librarian"}

	for _, roleName := range roles {
		var role models.Role
		err := db.FirstOrCreate(&role, models.Role{Name: roleName}).Error
		if err != nil {
			panic("Gagal melakukan seeding role: " + err.Error())
		}
	}
	log.Println("Seeding role berhasil!")
}
