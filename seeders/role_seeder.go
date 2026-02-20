package seeders

import (
	"belajar-go/config"
	"belajar-go/models"
	"log"
)

func SeedRoles() {
	roles := []string{"admin", "member", "librarian"}

	for _, roleName := range roles {
		var role models.Role
		err := config.DB.FirstOrCreate(&role, models.Role{Name: roleName}).Error
		if err != nil {
			panic("Gagal melakukan seeding role: " + err.Error())
		}
	}
	log.Println("Seeding role berhasil!")
}
