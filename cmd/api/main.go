package main

import (
	"net/http"
	"os"

	"demo-rrweb/internal/config"
	"demo-rrweb/internal/handler"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func main() {
	// 1. Inisialisasi Database
	config.ConnectDatabase()

	// 2. Setup Gin Router
	r := gin.Default()
	r.Use(CORSMiddleware())

	// 3. Daftarkan Endpoint (Sangat rapi karena memanggil dari folder handler)
	r.POST("/api/logs", handler.SaveLog)
	r.GET("/api/sessions", handler.GetSessionsList)    // API untuk list tabel Admin
	r.GET("/api/sessions/:id", handler.GetSessionByID) // API untuk player video
	r.DELETE("/api/sessions/:id", handler.DeleteSession)

	// 4. Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
