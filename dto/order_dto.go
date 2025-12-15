package dto

import "time"

// === Requests ===
type CreateOrderItemRequest struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int64 `json:"quantity"   binding:"required,gt=0"`
}

type CreateOrderRequest struct {
	Latitude  float64                  `json:"latitude"  binding:"required"`
	Longitude float64                  `json:"longitude" binding:"required"`
	Products  []CreateOrderItemRequest `json:"products"  binding:"required,min=1"`
}

type UpdateOrderRequest struct {
	// Dùng pointer để biết trường nào được gửi lên (để update từng phần)
	OrderID       int64   `json:"order_id" binding:"required"`
	PaymentStatus *string `json:"payment_status"`
	OrderStatus   *string `json:"order_status"`
}

type ReceiveOrderRequest struct {
	OrderID int64 `json:"order_id" binding:"required"`
}

// === Responses ===

type OrderSummaryResponse struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	PaymentStatus string    `json:"payment_status"`
	OrderStatus   string    `json:"order_status"`
	TotalAmount   float64   `json:"total_amount"`
	Thumbnail     string    `json:"thumbnail"` // URL ảnh
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	// (Thêm các trường khác nếu cần)
}

type OrderDetailResponse struct {
	OrderSummaryResponse
	UserName    string                    `json:"user_name"`
	UserPhone   string                    `json:"user_phone"`
	ShipperInfo *ShipperInfoResponse      `json:"shipper_info,omitempty"`
	Items       []OrderItemDetailResponse `json:"items"`
}

type OrderItemDetailResponse struct {
	ProductID    int64   `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductImage string  `json:"product_image"`
	Quantity     int64   `json:"quantity"`
	Price        float64 `json:"price"`
	Subtotal     float64 `json:"subtotal"`
}

type ShipperInfoResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type OrderListResponse struct {
	Orders     []OrderSummaryResponse `json:"orders"`
	Pagination Pagination             `json:"pagination"`
}
