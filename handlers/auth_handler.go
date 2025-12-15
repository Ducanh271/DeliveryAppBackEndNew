package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"example.com/delivery-app/dto"
	"example.com/delivery-app/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var req dto.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Trim spaces from email
	req.Email = strings.TrimSpace(req.Email)

	user, err := h.authService.SignUp(&req)
	if err != nil {
		if errors.Is(err, service.ErrEmailInUse) {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign up"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully, please verify your email",
		"user": gin.H{
			"id":      user.ID,
			"name":    user.Name,
			"email":   user.Email,
			"phone":   user.Phone,
			"address": user.Address,
			"role":    user.Role,
		},
	})
}

func (h *AuthHandler) VerifyOTPHandler(c *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.VerifyOTP(req.Email, req.OTP)

	if err != nil {
		if errors.Is(err, service.ErrInvalidOTP) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
			return
		} else if errors.Is(err, service.ErrOTPExpired) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP has expired"})
			return
		} else if errors.Is(err, service.ErrEmailNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Trim spaces from email
	req.Email = strings.TrimSpace(req.Email)

	resp, err := h.authService.Login(&req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) || errors.Is(err, service.ErrInvalidPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		} else if errors.Is(err, service.ErrUserInactive) {
			c.JSON(http.StatusForbidden, gin.H{"error": "User is inactive, verify your email"})
			return
		} else if errors.Is(err, service.ErrUserBanned) {
			c.JSON(http.StatusForbidden, gin.H{"error": "User is banned"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		log.Printf("Failed to login: %v\n", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// === 1. FORGOT PASSWORD HANDLER ===
func (h *AuthHandler) ForgotPasswordHandler(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.ForgotPassword(&req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			// Bất kỳ lỗi nào khác (tạo OTP, gửi email, CSDL)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset OTP"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to your email"})
}

// === 2. VERIFY RESET OTP HANDLER ===
func (h *AuthHandler) VerifyResetOTPHandler(c *gin.Context) {
	var req dto.VerifyResetOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	tokenString, err := h.authService.VerifyResetOTP(&req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, service.ErrInvalidOTP) || errors.Is(err, service.ErrOTPExpired) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reset_token": tokenString,
		"expires_in":  300, // 5 phút
	})
}

// === 3. RESET PASSWORD HANDLER ===
func (h *AuthHandler) ResetPasswordHandler(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	err := h.authService.ResetPassword(&req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidResetToken) || errors.Is(err, service.ErrTokenPurposeMismatch) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else
		// else if errors.Is(err, service.ErrPasswordTooWeak) {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})}
		{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update password"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	res, err := h.authService.RefreshToken(&req)

	if err != nil {
		// Lỗi 401 (nghiệp vụ)
		if errors.Is(err, service.ErrInvalidRefreshToken) ||
			errors.Is(err, service.ErrRefreshTokenExpired) ||
			errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			// Lỗi 500 (hệ thống)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not process token refresh"})
		}
		return
	}

	// 4. Thành công (200)
	c.JSON(http.StatusOK, res)
}
func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.LogOutRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	err := h.authService.LogOut(req.RefreshToken)

	if err != nil {
		if errors.Is(err, service.ErrInvalidRefreshToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not process logout"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
