package handler

import (
	"net/http"

	"github.com/emirrasyad/loopaffi-backend/internal/repository"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentRepo *repository.PaymentRepository
}

func NewPaymentHandler(paymentRepo *repository.PaymentRepository) *PaymentHandler {
	return &PaymentHandler{paymentRepo: paymentRepo}
}

func (h *PaymentHandler) GetAllPayments(c *gin.Context) {
	payments, err := h.paymentRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data pembayaran",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   payments,
	})
}

func (h *PaymentHandler) MarkAsPaid(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID pembayaran diperlukan",
		})
		return
	}

	if err := h.paymentRepo.MarkAsPaid(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui status pembayaran: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Pembayaran berhasil ditandai lunas",
	})
}

func (h *PaymentHandler) GetMyPayments(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}
	payments, err := h.paymentRepo.FindByAffiliateID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil data pembayaran"})
		return
	}
	if payments == nil {
		payments = []entity.Payment{}
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": payments})
}
