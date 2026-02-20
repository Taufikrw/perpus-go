package config

import (
	"belajar-go/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	timezone := os.Getenv("DB_TIMEZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslmode, timezone)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Gagal connect ke database! ", err)
	}

	database.AutoMigrate(&models.User{})
	database.AutoMigrate(&models.Role{})
	database.AutoMigrate(&models.Loan{})
	database.AutoMigrate(&models.Fine{})
	database.AutoMigrate(&models.BookCategory{})
	database.AutoMigrate(&models.Book{})
	database.AutoMigrate(&models.BookItem{})
	database.AutoMigrate(&models.Member{})

	DB = database
}
