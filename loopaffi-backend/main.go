package main

import (
	"backend/config"
	"backend/models"
	"backend/routes"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Inisialisasi Koneksi Database PostgreSQL
	config.ConnectDatabase()

	// 2. Auto-Migrate Database (Genap 10 Tabel sesuai Class Diagram & ERD)
	// GORM akan otomatis membuat tabel jika belum ada atau memperbarui strukturnya
	err := config.DB.AutoMigrate(
		// Tabel Master / Referensi
		&models.Role{},
		&models.User{},
		&models.Product{},
		&models.PaymentMethod{},
		&models.CommissionSetting{},

		// Tabel Transaksional
		&models.Sale{},
		&models.SaleItem{},
		&models.Commission{},
		&models.Payment{},
		&models.Notification{},
	)

	if err != nil {
		fmt.Println("Gagal melakukan migrasi database:", err)
	} else {
		fmt.Println("Migrasi 10 tabel sistem LoopAffi berhasil dilakukan!")
	}

	// 3. Setup Framework Gin
	router := gin.Default()

	// 4. CORS Middleware — agar frontend Next.js bisa akses API
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 5. Daftarkan Routes API (Endpoint untuk Postman/Frontend)
	routes.SetupRoutes(router)

	// 6. Jalankan Server (Ambil port dari os.Getenv)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server backend LoopAffi siap digunakan di http://localhost:%s\n", port)
	router.Run(":" + port)
}
