package handlers

import (
	"errors"
	"example.com/delivery-app/dto"
	"example.com/delivery-app/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ReviewHandler struct {
	reviewService *service.ReviewService
}

func NewReviewHandler(reviewService *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

// === CREATE REVIEW ===
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDRaw.(int64)

	var req dto.CreateReviewRequest

	// 1. Bind Form Data (vì có upload ảnh nên dùng ShouldBind thay vì ShouldBindJSON)
	// DTO đã có binding:"required" nên Gin sẽ tự validate
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review data: " + err.Error()})
		return
	}

	// 2. Validate rate thủ công nếu cần (mặc dù binding min/max đã lo)
	if req.Rate < 1 || req.Rate > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rate must be between 1 and 5"})
		return
	}

	// 3. Lấy files (optional)
	form, _ := c.MultipartForm()
	files := form.File["images"] // key "images" từ form-data

	// 4. Gọi Service
	review, err := h.reviewService.CreateReview(c.Request.Context(), userID, &req, files)
	if err != nil {
		if errors.Is(err, service.ErrUserCannotReview) || errors.Is(err, service.ErrReviewExists) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Review created successfully",
		"review":  review,
	})
}

// === GET REVIEWS BY PRODUCT ===
func (h *ReviewHandler) GetReviewsByProduct(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	resp, err := h.reviewService.GetReviewsByProduct(productID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
