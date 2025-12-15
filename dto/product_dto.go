package dto

import (
	"example.com/delivery-app/models"
	"time"
)

// CreateProductRequest là DTO cho `form:"..."`
type CreateProductRequest struct {
	Name        string  `form:"name"        binding:"required"`
	Description string  `form:"description" binding:"required"`
	Price       float64 `form:"price"       binding:"required,gt=0"`
	QtyInitial  int64   `form:"qty_initial" binding:"required,gte=0"`
	QtySold     int64   `form:"qty_sold"    binding:"gte=0"`
}

// ProductImageResponse là một phần của ProductDetailResponse
type ProductImageResponse struct {
	ID     int64  `json:"id"`
	URL    string `json:"url"`
	IsMain bool   `json:"is_main"`
}

// ProductDetailResponse là DTO trả về cho GetProductByID
type ProductDetailResponse struct {
	ID          int64                  `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Price       float64                `json:"price"`
	QtyInitial  int64                  `json:"qty_initial"`
	QtySold     int64                  `json:"qty_sold"`
	CreatedAt   time.Time              `json:"created_at"`
	Images      []ProductImageResponse `json:"images"`
	AvgRate     float64                `json:"avg_rate"`
	ReviewCount int                    `json:"review_count"`
}

//	type ProductOfListResponse struct {
//		ID          int64     `json:"id"`
//		Name        string    `json:"name"`
//		Price       float64   `json:"price"`
//		QtyInitial  int64     `json:"qty_initial"`
//		QtySold     int64     `json:"qty_sold"`
//		CreatedAt   time.Time `json:"created_at"`
//		MainImage   string    `json:"main_image"`
//		AvgRate     float64   `json:"avg_rate"`
//		ReviewCount int       `json:"review_count"`
//	}
type ProductResponse struct {
	ID          int64                 `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Price       float64               `json:"price"`
	QtyInitial  int64                 `json:"qty_initial"`
	QtySold     int64                 `json:"qty_sold"`
	CreatedAt   time.Time             `json:"created_at"`
	Images      []models.ProductImage `json:"images"`
	AvgRate     float64               `json:"avg_rate"`
	ReviewCount int64                 `json:"review_count"`
}
type ProductListResponse struct {
	Products   []ProductResponse `json:"products"`
	Pagination Pagination        `json:"pagination"`
}
