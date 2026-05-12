package routes

import (
	"github.com/emirrasyad/loopaffi-backend/internal/handler"
	"github.com/emirrasyad/loopaffi-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine           *gin.Engine
	authHandler      *handler.AuthHandler
	userHandler      *handler.UserHandler
	notifHandler     *handler.NotificationHandler
	saleHandler      *handler.SaleHandler
	commissionHandler *handler.CommissionHandler
	paymentHandler   *handler.PaymentHandler
	reportHandler    *handler.ReportHandler
	jwtSecret        string
}

func NewRouter(
	engine *gin.Engine,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	notifHandler *handler.NotificationHandler,
	saleHandler *handler.SaleHandler,
	commissionHandler *handler.CommissionHandler,
	paymentHandler *handler.PaymentHandler,
	reportHandler *handler.ReportHandler,
	jwtSecret string,
) *Router {
	return &Router{
		engine:           engine,
		authHandler:      authHandler,
		userHandler:      userHandler,
		notifHandler:     notifHandler,
		saleHandler:      saleHandler,
		commissionHandler: commissionHandler,
		paymentHandler:   paymentHandler,
		reportHandler:    reportHandler,
		jwtSecret:        jwtSecret,
	}
}

func (r *Router) Setup() {
	// Health check
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "app": "LoopAffi Backend"})
	})

	api := r.engine.Group("/api/v1")

	// Public routes
	auth := api.Group("/auth")
	{
		auth.POST("/login", r.authHandler.Login)
		auth.POST("/register", r.authHandler.Register)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(r.jwtSecret))
	{
		// Common notifications route for both roles
		protected.GET("/notifications", r.notifHandler.GetMyNotifications)
		protected.PUT("/notifications/:id/read", r.notifHandler.MarkRead)

		// Admin only
		admin := protected.Group("/admin")
		admin.Use(middleware.RoleMiddleware("admin"))
		{
			admin.GET("/users", r.userHandler.GetAllUsers)
			admin.POST("/sales", r.saleHandler.CreateSale)
			admin.GET("/sales", r.saleHandler.GetAllSales)
			admin.GET("/commissions", r.commissionHandler.GetAllCommissions)
			admin.GET("/payments", r.paymentHandler.GetAllPayments)
			admin.PUT("/payments/:id/pay", r.paymentHandler.MarkAsPaid)
			admin.GET("/reports", r.reportHandler.GetReport)
		}

		// Affiliate only
		affiliate := protected.Group("/affiliate")
		affiliate.Use(middleware.RoleMiddleware("affiliate"))
		{
			affiliate.GET("/sales", r.saleHandler.GetMySales)
			affiliate.GET("/commissions", r.commissionHandler.GetMyCommissions)
			affiliate.GET("/payments", r.paymentHandler.GetMyPayments)
		}
	}
}
