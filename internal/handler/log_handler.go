package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"demo-rrweb/internal/config"
	"demo-rrweb/internal/model"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Simpan data dari Vue ke Postgres
func SaveLog(c *gin.Context) {
	// Ambil context dari request untuk tracing
	ctx := c.Request.Context()

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		// CONTOH LOG ERROR LENGKAP
		slog.Error("Gagal membaca input JSON dari frontend",
			slog.String("error", err.Error()),
			slog.String("trace_id", trace.SpanFromContext(ctx).SpanContext().TraceID().String()),
		)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ngitung kecepatan konversi JSON
	tracer := otel.Tracer("rrweb-handler")
	spanCtx, span := tracer.Start(ctx, "Proses Format JSON")

	// Ubah array object menjadi string mentah untuk masuk ke kolom JSONB
	eventsRaw, _ := json.Marshal(input["events"])

	logBaru := model.SessionLog{
		Events:    string(eventsRaw),
		CreatedAt: time.Now(),
	}

	span.End() // Akhiri span untuk proses konversi JSON

	if err := config.DB.WithContext(spanCtx).Create(&logBaru).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan ke database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil disimpan!"})
}

// Mengambil Daftar Sesi (Hanya ID dan Waktu untuk Tabel Admin)
func GetSessionsList(c *gin.Context) {
	ctx := c.Request.Context()

	// Kita buat struct kecil sementara agar database tidak perlu
	// menarik data JSON (events) yang ukurannya bisa ber-megabyte
	type SessionSummary struct {
		ID        uint      `json:"id"`
		CreatedAt time.Time `json:"created_at"`
	}

	var summaries []SessionSummary

	// Select membatasi kolom yang ditarik, membuat query super cepat
	if err := config.DB.WithContext(ctx).Model(&model.SessionLog{}).Select("id", "created_at").Order("created_at desc").Find(&summaries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, summaries)
}

// Mengambil Data JSON Spesifik berdasarkan ID (Untuk Replayer)
func GetSessionByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id") // Mengambil ID dari URL (/api/sessions/1)
	var session model.SessionLog

	if err := config.DB.WithContext(ctx).First(&session, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sesi tidak ditemukan"})
		return
	}

	// Ubah kembali string JSONB dari DB menjadi array asli
	var eventsData []map[string]interface{}
	json.Unmarshal([]byte(session.Events), &eventsData)

	c.JSON(http.StatusOK, eventsData)
}

// Menghapus Sesi dari Database berdasarkan ID
func DeleteSession(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id") // Tangkap ID dari URL

	// Perintah GORM untuk menghapus baris di PostgreSQL berdasarkan ID
	if err := config.DB.WithContext(ctx).Delete(&model.SessionLog{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data dari database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data rekaman berhasil dihapus!"})
}
