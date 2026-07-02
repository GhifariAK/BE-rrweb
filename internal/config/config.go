package config

import (
	"fmt"
	"log"
	"os"

	"demo-rrweb/internal/model"

	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

// DB adalah variabel global untuk dipakai oleh handler nanti
var DB *gorm.DB

// GetEnv adalah helper untuk mengambil nilai dari .env dengan nilai cadangan (fallback)
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		GetEnv("DB_HOST", "localhost"), GetEnv("DB_USER", "postgres"), GetEnv("DB_PASS", ""),
		GetEnv("DB_NAME", "rrweb_db"), GetEnv("DB_PORT", "5432"),
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database PostgreSQL: ", err)
	}

	// Pasang pengawas di semua query database
	if err := database.Use(tracing.NewPlugin()); err != nil {
		log.Fatal("Gagal memasang sensor OTel GORM:", err)
	}

	// Otomatis membuat/update tabel dari struct SessionLog
	database.AutoMigrate(&model.SessionLog{})

	DB = database
	fmt.Println("Database PostgreSQL berhasil terhubung dan termigrasi!")
}
