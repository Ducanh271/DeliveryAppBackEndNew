package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"example.com/delivery-app/dto"
	"example.com/delivery-app/models"
	"example.com/delivery-app/repository"
)

type OrderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	userRepo    repository.UserRepository
	uow         repository.UnitOfWork
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	userRepo repository.UserRepository,
	uow repository.UnitOfWork,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
		uow:         uow,
	}
}

// === 1. CREATE ORDER ===
func (s *OrderService) CreateOrder(userID int64, req *dto.CreateOrderRequest) error {
	var orderItems []models.OrderItem
	var totalAmount float64
	var thumbnailID int

	// 1. Validate sản phẩm và tính tiền
	for i, itemReq := range req.Products {

		product, err := s.productRepo.GetByID(itemReq.ProductID)
		currentStock := product.QtyInitial - product.QtySold
		if currentStock < itemReq.Quantity {
			return fmt.Errorf("sản phẩm '%s' không đủ hàng (còn lại: %d)", product.Name, currentStock)
		}
		if err != nil {
			return fmt.Errorf("product %d not found", itemReq.ProductID)
		}

		// Logic lấy Thumbnail (Lấy ảnh main của sản phẩm đầu tiên)
		if i == 0 {
			images, _ := s.productRepo.GetImagesByProductID(product.ID)
			for _, img := range images {
				if img.IsMain {
					thumbnailID = int(img.ID)
					break
				}
			}
			// Fallback: nếu không có main, lấy ảnh đầu
			if thumbnailID == 0 && len(images) > 0 {
				thumbnailID = int(images[0].ID)
			}
		}

		totalAmount += product.Price * float64(itemReq.Quantity)
		orderItems = append(orderItems, models.OrderItem{
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			Price:     product.Price, // Lưu giá tại thời điểm mua
		})
	}

	// 2. Tạo struct Order
	order := &models.Order{
		UserID:        userID,
		PaymentStatus: "unpaid",
		OrderStatus:   "pending",
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		TotalAmount:   totalAmount,
		ThumbnailID:   thumbnailID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// 3. Transaction (UoW)
	err := s.uow.Execute(func(repoProvider func(any) any) error {
		ordRepo := repoProvider((*repository.OrderRepository)(nil)).(repository.OrderRepository)
		prodRepo := repoProvider((*repository.ProductRepository)(nil)).(repository.ProductRepository) // [LƯU Ý] Cần thêm dòng này
		// Lưu Order
		orderID, err := ordRepo.CreateOrder(order)
		if err != nil {
			return err
		}

		// Lưu Order Items
		for _, item := range orderItems {
			item.OrderID = orderID
			if err := ordRepo.CreateOrderItem(&item); err != nil {
				return err
			}
			if err := prodRepo.UpdateQtySold(item.ProductID, item.Quantity); err != nil {
				return fmt.Errorf("lỗi cập nhật kho cho sp %d: %w", item.ProductID, err)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	return nil
}

// === 2. GET ORDER DETAIL ===
func (s *OrderService) GetOrderDetail(orderID int64, userID int64, role string) (*dto.OrderDetailResponse, error) {
	// 1. Check quyền (nếu là customer)
	if role == "customer" {
		isOwner, _ := s.orderRepo.CheckOrderOwnership(userID, orderID)
		if !isOwner {
			return nil, ErrNotYourOrder
		}
	}
	// (Nếu là shipper, có thể check thêm xem có phải đơn mình nhận không)
	if role == "shipper" {
		isOwner, _ := s.orderRepo.CheckShipperOwnership(userID, orderID)
		if !isOwner {
			return nil, ErrNotYourOrder
		}
	}
	// 2. Lấy Order
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	// 3. Lấy Items
	items, err := s.orderRepo.GetItemsByOrderID(orderID)
	if err != nil {
		return nil, err
	}

	// 4. Mapping Items (kèm tên và ảnh sản phẩm)
	var itemResponses []dto.OrderItemDetailResponse
	for _, item := range items {
		product, _ := s.productRepo.GetByID(item.ProductID)

		// Lấy ảnh đại diện cho item
		images, _ := s.productRepo.GetImagesByProductID(item.ProductID)
		imgUrl := ""
		if len(images) > 0 {
			imgUrl = images[0].URL
		}

		productName := "Unknown Product"
		if product != nil {
			productName = product.Name
		}

		itemResponses = append(itemResponses, dto.OrderItemDetailResponse{
			ProductID:    item.ProductID,
			ProductName:  productName,
			ProductImage: imgUrl,
			Quantity:     item.Quantity,
			Price:        item.Price,
			Subtotal:     item.Price * float64(item.Quantity),
		})
	}

	// 5. Lấy User Info & Shipper Info
	user, _ := s.userRepo.GetUserByID(order.UserID)
	shipperID, shipperName, shipperPhone, _ := s.orderRepo.GetShipperInfoByOrderID(orderID)

	// 6. Final Response
	resp := &dto.OrderDetailResponse{
		OrderSummaryResponse: dto.OrderSummaryResponse{
			ID:            order.ID,
			UserID:        order.UserID,
			PaymentStatus: order.PaymentStatus,
			OrderStatus:   order.OrderStatus,
			TotalAmount:   order.TotalAmount,
			CreatedAt:     order.CreatedAt,
			UpdatedAt:     order.UpdatedAt,
			Longitude:     order.Longitude,
			Latitude:      order.Latitude,
		},
		UserName:  user.Name,
		UserPhone: user.Phone,
		Items:     itemResponses,
	}

	if shipperID != 0 {
		resp.ShipperInfo = &dto.ShipperInfoResponse{
			ID:    shipperID,
			Name:  shipperName,
			Phone: shipperPhone,
		}
	}

	return resp, nil
}

// === 3. LIST ORDERS (Generic) ===
func (s *OrderService) GetListOrders(filter map[string]interface{}, page, limit int) (*dto.OrderListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	orders, total, err := s.orderRepo.GetOrders(filter, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	ThumbnailIDs := make([]int64, 0)
	for _, o := range orders {
		if o.ThumbnailID != 0 {
			ThumbnailIDs = append(ThumbnailIDs, int64(o.ThumbnailID))
		}
	}
	thumbnailMap, err := s.orderRepo.GetImageURLsByIDs(ThumbnailIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get thumbnail URLs: %w", err)
	}

	// Mapping to DTO
	var orderSummaries []dto.OrderSummaryResponse
	for _, o := range orders {
		// Tuy nhiên để tối ưu, ta có thể bỏ qua hoặc dùng Batch Query như bên Product)

		orderSummaries = append(orderSummaries, dto.OrderSummaryResponse{
			ID:            o.ID,
			UserID:        o.UserID,
			PaymentStatus: o.PaymentStatus,
			OrderStatus:   o.OrderStatus,
			TotalAmount:   o.TotalAmount,
			Thumbnail:     thumbnailMap[int64(o.ThumbnailID)],
			UpdatedAt:     o.UpdatedAt,
			Latitude:      o.Latitude,
			Longitude:     o.Longitude,
		})
	}

	totalPages := (total + limit - 1) / limit
	return &dto.OrderListResponse{
		Orders: orderSummaries,
		Pagination: dto.Pagination{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	}, nil
}

// === 4. SHIPPER: RECEIVE ORDER ===
func (s *OrderService) ReceiveOrder(shipperID int64, orderID int64) error {
	// 1. Check đơn có tồn tại và đang 'processing' không
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order.OrderStatus != "processing" {
		return ErrOrderNotProcessing
	}

	// 2. Check Shipper đã nhận quá 10 đơn chưa
	count, err := s.orderRepo.CountActiveOrdersByShipper(shipperID)
	if err != nil {
		return err
	}
	if count >= 10 {
		return ErrMaxOrdersReached
	}

	// 3. Update
	return s.orderRepo.UpdateShipperForOrder(orderID, shipperID)
}

// === 5. SHIPPER: UPDATE STATUS ===
func (s *OrderService) UpdateOrder(shipperID int64, req *dto.UpdateOrderRequest) error {
	// 1. Check quyền (Shipper phải là người đang ship đơn này)
	isOwner, err := s.orderRepo.CheckShipperOwnership(shipperID, req.OrderID)
	if err != nil {
		return err
	}
	if !isOwner {
		return ErrOrderNotOwned
	}

	// 2. Build map update
	updates := make(map[string]interface{})
	if req.PaymentStatus != nil {
		updates["payment_status"] = *req.PaymentStatus
	}
	if req.OrderStatus != nil {
		updates["order_status"] = *req.OrderStatus
	}

	if len(updates) == 0 {
		return nil
	}

	return s.orderRepo.UpdateOrder(req.OrderID, updates)
}

// === 6. ADMIN: ACCEPT ORDER (Pending -> Processing) ===
func (s *OrderService) AdminAcceptOrder(orderID int64) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order.OrderStatus != "pending" {
		return ErrOrderNotPending
	}
	// Check order tồn tại...
	return s.orderRepo.UpdateOrder(orderID, map[string]interface{}{
		"order_status": "processing",
	})
}

// === 7. CUSTOMER: CANCEL ORDER ===
func (s *OrderService) CancelOrder(userID int64, orderID int64) error {
	// 1. Check quyền
	isOwner, _ := s.orderRepo.CheckOrderOwnership(userID, orderID)
	if !isOwner {
		return ErrNotYourOrder
	}

	// 2. Check trạng thái
	order, _ := s.orderRepo.GetByID(orderID)
	if order.OrderStatus != "pending" {
		return ErrCannotCancel
	}

	// 3. Cancel
	return s.orderRepo.UpdateOrder(orderID, map[string]interface{}{
		"order_status": "cancelled",
	})
}

// === 8. DASHBOARD STATS ===
func (s *OrderService) GetStats() (int64, float64, error) {
	return s.orderRepo.GetStats()
}
