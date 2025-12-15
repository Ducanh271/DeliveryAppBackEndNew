package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"example.com/delivery-app/dto"
	"example.com/delivery-app/service"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *service.ProductService
}

// NewProductHandler là hàm khởi tạo
func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// === Create Product ===
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest

	// 1. Bind Form Data (các trường text)
	// Binding sẽ tự động check 'required', 'gt=0'...
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data: " + err.Error()})
		return
	}

	// 2. Lấy Files
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can't read form data for images"})
		return
	}
	files := form.File["images"] // key là "images"

	// 3. Gọi Service
	product, urls, err := h.productService.CreateProduct(c.Request.Context(), &req, files)

	// 4. Xử lý lỗi từ Service
	if err != nil {
		if errors.Is(err, service.ErrNoImages) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			// Các lỗi khác (upload, db...) coi là lỗi server
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create product: " + err.Error()})
		}
		return
	}

	// 5. Thành công
	c.JSON(http.StatusCreated, gin.H{
		"message": "Created new product successfully",
		"product": product,
		"images":  urls,
	})
}

// === Get Product By ID ===
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Gọi Service
	productRes, err := h.productService.GetProductDetail(id)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product detail"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": productRes})
}

// === Delete Product ===
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Gọi Service
	err = h.productService.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found to delete"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// === Get Products (Bao gồm cả Search và GetAll) ===
func (h *ProductHandler) GetProducts(c *gin.Context) {
	// Lấy params từ URL
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	keyword := c.Query("q") // Nếu có q thì là search, không thì là get all

	// Gọi Service (Service đã xử lý logic search/getall bên trong)
	response, err := h.productService.GetProducts(keyword, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// === Get Total Products (API phụ) ===
func (h *ProductHandler) GetTotalProducts(c *gin.Context) {
	count, err := h.productService.GetTotalProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get count"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"number of products": count})
}
