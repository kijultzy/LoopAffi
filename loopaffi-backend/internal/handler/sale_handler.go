package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/emirrasyad/loopaffi-backend/internal/entity"
	"github.com/emirrasyad/loopaffi-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SaleHandler struct {
	saleRepo       *repository.SaleRepository
	notifRepo      *repository.NotificationRepository
	commissionRepo *repository.CommissionRepository
	paymentRepo    *repository.PaymentRepository
}

func NewSaleHandler(
	saleRepo *repository.SaleRepository,
	notifRepo *repository.NotificationRepository,
	commissionRepo *repository.CommissionRepository,
	paymentRepo *repository.PaymentRepository,
) *SaleHandler {
	return &SaleHandler{
		saleRepo:       saleRepo,
		notifRepo:      notifRepo,
		commissionRepo: commissionRepo,
		paymentRepo:    paymentRepo,
	}
}

func (h *SaleHandler) CreateSale(c *gin.Context) {
	var input struct {
		Amount      float64 `json:"amount" binding:"required"`
		AffiliateID string  `json:"affiliate_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Input tidak valid"})
		return
	}

	now := time.Now()
	dateStr := now.Format("2006-01-02")

	sale := &entity.Sale{
		ID:          uuid.New().String(),
		Date:        dateStr,
		Amount:      input.Amount,
		AffiliateID: input.AffiliateID,
		Status:      "completed",
	}

	if err := h.saleRepo.Create(sale); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mencatat penjualan"})
		return
	}

	// Hitung komisi (10%)
	commissionAmount := input.Amount * 0.1

	// Buat commission record di database
	commission := &entity.Commission{
		ID:          "COM-" + uuid.New().String()[:8],
		SaleID:      sale.ID,
		AffiliateID: input.AffiliateID,
		Amount:      commissionAmount,
		Date:        dateStr,
	}
	if err := h.commissionRepo.Create(commission); err != nil {
		fmt.Printf("Gagal membuat komisi: %v\n", err)
	}

	// Buat payment record di database (status pending)
	payment := &entity.Payment{
		ID:          "PAY-" + uuid.New().String()[:8],
		AffiliateID: input.AffiliateID,
		Amount:      commissionAmount,
		Date:        dateStr,
		Status:      "pending",
	}
	if err := h.paymentRepo.Create(payment); err != nil {
		fmt.Printf("Gagal membuat pembayaran: %v\n", err)
	}

	// Buat notifikasi untuk affiliator
	message := fmt.Sprintf("Penjualan baru tercatat! Anda mendapatkan komisi sebesar Rp %.0f", commissionAmount)

	notif := &entity.Notification{
		UserID:  input.AffiliateID,
		Message: message,
	}

	if err := h.notifRepo.Create(notif); err != nil {
		// Log error tapi jangan gagalkan respon penjualan
		fmt.Printf("Gagal membuat notifikasi: %v\n", err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Penjualan berhasil dicatat",
		"data":    sale,
	})
}

func (h *SaleHandler) GetAllSales(c *gin.Context) {
	sales, err := h.saleRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil data penjualan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   sales,
	})
}

func (h *SaleHandler) GetMySales(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}
	sales, err := h.saleRepo.FindByAffiliateID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil data penjualan"})
		return
	}
	if sales == nil {
		sales = []entity.Sale{}
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": sales})
}
