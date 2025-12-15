package dto

import "time"

// Request tạo review (Form data để upload ảnh)
type CreateReviewRequest struct {
	ProductID int64  `form:"product_id" binding:"required"`
	OrderID   int64  `form:"order_id"   binding:"required"`
	Rate      int8   `form:"rate"       binding:"required,min=1,max=5"`
	Content   string `form:"content"    binding:"required"`
}

// Response hiển thị ảnh
type ReviewImageResponse struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

// Response hiển thị 1 review đầy đủ
type ReviewResponse struct {
	ID        int64                 `json:"id"`
	UserID    int64                 `json:"user_id"`
	UserName  string                `json:"user_name"`
	Rate      int8                  `json:"rate"`
	Content   string                `json:"content"`
	CreatedAt time.Time             `json:"created_at"`
	Images    []ReviewImageResponse `json:"images"`
}

// Response danh sách review (phân trang)
type ReviewListResponse struct {
	Reviews    []ReviewResponse `json:"reviews"`
	Pagination Pagination       `json:"pagination"`
}
