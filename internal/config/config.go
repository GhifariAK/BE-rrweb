package config

import (
	"fmt"
	"log"
	"os"

	"demo-rrweb/internal/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB adalah variabel global untuk dipakai oleh handler nanti
var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database PostgreSQL: ", err)
	}

	// Otomatis membuat/update tabel dari struct SessionLog
	database.AutoMigrate(&model.SessionLog{})

	DB = database
	fmt.Println("Database PostgreSQL berhasil terhubung dan termigrasi!")
}
