package models

import "time"

type Review struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	UserID    int64     `json:"user_id"`
	OrderID   int64     `json:"order_id"`
	Rate      int8      `json:"rate"` // 1-5
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// Struct phụ để hứng dữ liệu join
type ReviewDetail struct {
	Review
	UserName string `json:"user_name"` // Join từ bảng Users
}

type ReviewImage struct {
	ID       int64  `json:"id"`
	ReviewID int64  `json:"review_id"`
	URL      string `json:"url"`
	PublicID string `json:"public_id"`
}
