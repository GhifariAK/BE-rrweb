package main

import (
	"context"
	"log"
	"net/http"

	"demo-rrweb/internal/config"
	"demo-rrweb/internal/handler"
	"demo-rrweb/internal/telemetry"

	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin" //middlewaree khusus Gin dari Otel
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func main() {
	// 1. Init Opentelemetry
	tel, err := telemetry.InitTelemetry()
	if err != nil {
		log.Fatal("Gagal menyalakan OTel:", err)
	}

	defer func() {
		ctx := context.Background()
		tel.Trace(ctx)
		tel.Log(ctx)
	}()

	// 2. Init Logger
	config.InitLogger()

	// 3. Init Database
	config.ConnectDatabase()

	// 4. Setup Gin Router
	r := gin.Default()
	r.Use(CORSMiddleware())

	// 5. Daftarkan Middleware Otel untuk Gin
	r.Use(otelgin.Middleware("rrweb-backend"))

	// 6. Daftarkan Endpoint (Sangat rapi karena memanggil dari folder handler)
	r.POST("/api/logs", handler.SaveLog)
	r.GET("/api/sessions", handler.GetSessionsList)    // API untuk list tabel Admin
	r.GET("/api/sessions/:id", handler.GetSessionByID) // API untuk player video
	r.DELETE("/api/sessions/:id", handler.DeleteSession)

	// 7. Jalankan Server
	port := config.GetEnv("PORT", "8080")
	r.Run(":" + port)
}
