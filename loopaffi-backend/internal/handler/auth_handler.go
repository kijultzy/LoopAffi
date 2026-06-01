package handler

import (
	"github.com/emirrasyad/loopaffi-backend/internal/dto"
	"github.com/emirrasyad/loopaffi-backend/internal/response"
	"github.com/emirrasyad/loopaffi-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login godoc
// @Summary Login user
// @Description Login dengan email dan password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse
// @Failure 401 {object} response.StandardResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Request tidak valid: "+err.Error())
		return
	}

	result, err := h.authService.Login(req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.OK(c, "Login berhasil", result)
}

// Register godoc
// @Summary Register affiliator
// @Description Daftar sebagai affiliator baru
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register data"
// @Success 201 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse
// @Failure 409 {object} response.StandardResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Request tidak valid: "+err.Error())
		return
	}

	result, err := h.authService.Register(req)
	if err != nil {
		// Email duplicate = conflict
		if err.Error() == "email sudah terdaftar" {
			response.Conflict(c, err.Error())
			return
		}
		response.InternalServerError(c, "Gagal mendaftarkan akun: "+err.Error())
		return
	}

	response.Created(c, "Registrasi berhasil", result)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Request tidak valid: "+err.Error())
		return
	}

	err := h.authService.ForgotPassword(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "Password berhasil diubah. Silakan login dengan password baru.", nil)
}
