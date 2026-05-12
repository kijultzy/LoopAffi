package main

import (
	"fmt"
	"log"
	"os"

	"github.com/emirrasyad/loopaffi-backend/internal/config"
	"github.com/emirrasyad/loopaffi-backend/internal/database"
	"github.com/emirrasyad/loopaffi-backend/internal/handler"
	"github.com/emirrasyad/loopaffi-backend/internal/repository"
	"github.com/emirrasyad/loopaffi-backend/internal/routes"
	"github.com/emirrasyad/loopaffi-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Connect to Database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 3. Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	notifRepo := repository.NewNotificationRepository(db)
	saleRepo := repository.NewSaleRepository(db)
	commissionRepo := repository.NewCommissionRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	// 4. Initialize Services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)

	// 5. Initialize Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userRepo)
	notifHandler := handler.NewNotificationHandler(notifRepo)
	saleHandler := handler.NewSaleHandler(saleRepo, notifRepo, commissionRepo, paymentRepo)
	commissionHandler := handler.NewCommissionHandler(commissionRepo)
	paymentHandler := handler.NewPaymentHandler(paymentRepo)
	reportHandler := handler.NewReportHandler(userRepo, saleRepo, commissionRepo, paymentRepo)

	// 6. Setup Gin Router
	engine := gin.Default()

	// CORS Middleware
	engine.Use(func(c *gin.Context) {
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

	api := engine.Group("/api")
	api.GET("/migrate", func(c *gin.Context) {
		sqlBytes, err := os.ReadFile("migrations/001_create_users.up.sql")
		if err != nil {
			c.JSON(500, gin.H{"error": "Cannot read migration file: " + err.Error()})
			return
		}
		_, err = db.Exec(string(sqlBytes))
		if err != nil {
			c.JSON(500, gin.H{"error": "Migration failed: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Database migrated successfully!"})
	})

	// 7. Setup Routes
	router := routes.NewRouter(
		engine,
		authHandler,
		userHandler,
		notifHandler,
		saleHandler,
		commissionHandler,
		paymentHandler,
		reportHandler,
		cfg.JWTSecret,
	)
	router.Setup()

	// 8. Run Server
	port := cfg.AppPort
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server backend LoopAffi siap digunakan di port %s\n", port)
	if err := engine.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
