package handler

import (
	"net/http"

	"github.com/emirrasyad/loopaffi-backend/internal/repository"
	"github.com/gin-gonic/gin"
)

// ReportAffiliateRow represents a single affiliate's aggregated report data
type ReportAffiliateRow struct {
	AffiliateID       string  `json:"affiliate_id"`
	AffiliateName     string  `json:"affiliate_name"`
	AffiliateEmail    string  `json:"affiliate_email"`
	TotalSales        float64 `json:"total_sales"`
	SalesCount        int     `json:"sales_count"`
	TotalCommission   float64 `json:"total_commission"`
	PaidCommission    float64 `json:"paid_commission"`
	PendingCommission float64 `json:"pending_commission"`
}

type ReportHandler struct {
	userRepo       *repository.UserRepository
	saleRepo       *repository.SaleRepository
	commissionRepo *repository.CommissionRepository
	paymentRepo    *repository.PaymentRepository
}

func NewReportHandler(
	userRepo *repository.UserRepository,
	saleRepo *repository.SaleRepository,
	commissionRepo *repository.CommissionRepository,
	paymentRepo *repository.PaymentRepository,
) *ReportHandler {
	return &ReportHandler{
		userRepo:       userRepo,
		saleRepo:       saleRepo,
		commissionRepo: commissionRepo,
		paymentRepo:    paymentRepo,
	}
}

func (h *ReportHandler) GetReport(c *gin.Context) {
	// Get all affiliate users
	users, err := h.userRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil data user"})
		return
	}

	// Get all sales, commissions, payments
	sales, err := h.saleRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil data penjualan"})
		return
	}

	commissions, err := h.commissionRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil data komisi"})
		return
	}

	payments, err := h.paymentRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil data pembayaran"})
		return
	}

	// Build report rows only for affiliate users
	var report []ReportAffiliateRow
	for _, user := range users {
		if user.RoleID != "affiliate" {
			continue
		}

		row := ReportAffiliateRow{
			AffiliateID:    user.ID,
			AffiliateName:  user.Name,
			AffiliateEmail: user.Email,
		}

		// Aggregate sales
		for _, s := range sales {
			if s.AffiliateID == user.ID {
				row.TotalSales += s.Amount
				row.SalesCount++
			}
		}

		// Aggregate commissions
		for _, c := range commissions {
			if c.AffiliateID == user.ID {
				row.TotalCommission += c.Amount
			}
		}

		// Aggregate payments
		for _, p := range payments {
			if p.AffiliateID == user.ID {
				if p.Status == "paid" {
					row.PaidCommission += p.Amount
				} else if p.Status == "pending" {
					row.PendingCommission += p.Amount
				}
			}
		}

		report = append(report, row)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   report,
	})
}
