package config

import (
	"fmt"
	"log/slog"
	"os"

	"demo-rrweb/internal/model"

	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

// DB adalah variabel global untuk dipakai oleh handler nanti
var DB *gorm.DB

// Mode json untuk Log di terminal
func InitLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger) // Jadikan default untuk seluruh aplikasi
}

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
		slog.Warn("Peringatan: File .env tidak ditemukan")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		GetEnv("DB_HOST", "localhost"), GetEnv("DB_USER", "postgres"), GetEnv("DB_PASS", ""),
		GetEnv("DB_NAME", "rrweb_db"), GetEnv("DB_PORT", "5432"),
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Gagal terhubung ke database PostgreSQL: ", slog.String("error", err.Error()))
	}

	// Pasang pengawas di semua query database
	if err := database.Use(tracing.NewPlugin()); err != nil {
		slog.Error("Gagal memasang sensor OTel GORM:", slog.String("error", err.Error()))
	}

	// Otomatis membuat/update tabel dari struct SessionLog
	database.AutoMigrate(&model.SessionLog{})

	DB = database
	slog.Info("Database PostgreSQL berhasil terhubung dan termigrasi!")
}
