package handlers

import (
	"errors"
	"example.com/delivery-app/dto"
	"example.com/delivery-app/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// UserHandler xử lý các request quản lý user
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler là hàm khởi tạo
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// === TỪ CreateShipper ===
func (h *UserHandler) CreateShipper(c *gin.Context) {
	var req dto.CreateShipperRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("CreateShipper request:", req)
		return
	}

	user, err := h.userService.CreateShipper(&req)
	if err != nil {
		if errors.Is(err, service.ErrEmailInUse) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create shipper"})
		}
		log.Println("CreateShipper error (01):", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Shipper created successfully", "id": user.ID})
}

// === TỪ ProfileHandler ===
func (h *UserHandler) Profile(c *gin.Context) {
	// Lấy userID từ middleware (đã xác thực)
	userIDraw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDraw.(int64) // Ép kiểu

	userProfile, err := h.userService.GetUserProfile(userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get profile"})
		}
		return
	}

	c.JSON(http.StatusOK, userProfile)
}

// === TỪ BanUserAccountHandler ===
func (h *UserHandler) BanUser(c *gin.Context) {
	userIDstr := c.Param("id")
	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	log.Println("BanUser userID:", userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.userService.BanUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to ban user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User banned successfully"})
}

// === TỪ UnBanUserAccountHandler ===
func (h *UserHandler) UnbanUser(c *gin.Context) {
	userIDstr := c.Param("id")
	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.userService.UnbanUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unban user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User unbanned successfully"})
}

// === TỪ GetAllCustomersHandler ===
func (h *UserHandler) GetAllCustomers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	response, err := h.userService.GetUsersByRole("customer", page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get customers"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// === TỪ GetAllShippersHandler ===
func (h *UserHandler) GetAllShippers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	response, err := h.userService.GetUsersByRole("shipper", page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get shippers"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// === TỪ GetNumberOf...Handler ===
// (Gộp thành 1 hàm dashboard)
func (h *UserHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.userService.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
