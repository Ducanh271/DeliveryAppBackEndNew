package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"example.com/delivery-app/models"
)

type OrderRepository interface {
	// Write (Dùng trong UoW)
	CreateOrder(order *models.Order) (int64, error)
	CreateOrderItem(item *models.OrderItem) error
	UpdateOrder(orderID int64, updates map[string]interface{}) error
	UpdateShipperForOrder(orderID int64, shipperID int64) error

	// Read
	GetByID(id int64) (*models.Order, error)
	GetItemsByOrderID(orderID int64) ([]models.OrderItem, error)

	// Check
	CheckOrderOwnership(userID, orderID int64) (bool, error)
	CheckShipperOwnership(shipperID, orderID int64) (bool, error)
	CountActiveOrdersByShipper(shipperID int64) (int, error)

	// List & Filter (Hàm đa năng thay thế cho GetAllOrders, GetOrdersByShipper...)
	GetImageURLsByIDs(imageIDs []int64) (map[int64]string, error)
	GetOrders(filter map[string]interface{}, page, limit int) ([]models.Order, int, error)

	// Stats & Info
	GetShipperInfoByOrderID(orderID int64) (int64, string, string, error)
	GetStats() (int64, float64, error)

	WithTx(tx *sql.Tx) OrderRepository
}

type orderRepo struct {
	db DBTX
}

func NewOrderRepository(db DBTX) OrderRepository {
	return &orderRepo{db: db}
}

func (r *orderRepo) WithTx(tx *sql.Tx) OrderRepository {
	return &orderRepo{db: tx}
}

// --- IMPLEMENTATION ---

// 1. CREATE

func (r *orderRepo) CreateOrder(o *models.Order) (int64, error) {
	query := `INSERT INTO orders (user_id, payment_status, order_status, latitude, longitude, total_amount, thumbnail_id, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := r.db.Exec(query, o.UserID, o.PaymentStatus, o.OrderStatus, o.Latitude, o.Longitude, o.TotalAmount, o.ThumbnailID, o.CreatedAt, o.UpdatedAt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *orderRepo) CreateOrderItem(item *models.OrderItem) error {
	query := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, item.OrderID, item.ProductID, item.Quantity, item.Price)
	return err
}

// 2. UPDATE

// UpdateOrder là hàm update động, cực kỳ mạnh mẽ
func (r *orderRepo) UpdateOrder(orderID int64, updates map[string]interface{}) error {
	query := "UPDATE orders SET "
	var args []interface{}
	var i int

	for key, val := range updates {
		if i > 0 {
			query += ", "
		}
		query += key + " = ?"
		args = append(args, val)
		i++
	}

	// Luôn cập nhật updated_at
	query += ", updated_at = ? WHERE id = ?"
	args = append(args, time.Now(), orderID)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *orderRepo) UpdateShipperForOrder(orderID int64, shipperID int64) error {
	// Logic: Chỉ nhận đơn đang 'processing' (logic này có thể để ở Service check cũng được, nhưng để ở đây cho chắc chắn)
	query := "UPDATE orders SET shipper_id = ?, order_status = 'shipping', updated_at = ? WHERE id = ? AND order_status = 'processing'"
	res, err := r.db.Exec(query, shipperID, time.Now(), orderID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows // Không tìm thấy đơn phù hợp để nhận
	}
	return nil
}

// 3. READ SINGLE

func (r *orderRepo) GetByID(id int64) (*models.Order, error) {
	query := `SELECT id, user_id, shipper_id, payment_status, order_status, latitude, longitude, total_amount, thumbnail_id, created_at, updated_at 
			  FROM orders WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var o models.Order
	// Xử lý shipper_id có thể NULL
	var shipperID sql.NullInt64

	err := row.Scan(&o.ID, &o.UserID, &shipperID, &o.PaymentStatus, &o.OrderStatus, &o.Latitude, &o.Longitude, &o.TotalAmount, &o.ThumbnailID, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if shipperID.Valid {
		o.ShipperID = shipperID.Int64
	}

	return &o, nil
}

func (r *orderRepo) GetItemsByOrderID(orderID int64) ([]models.OrderItem, error) {
	query := `SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = ?`
	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// 4. CHECK

func (r *orderRepo) CheckOrderOwnership(userID, orderID int64) (bool, error) {
	var exists int
	err := r.db.QueryRow("SELECT 1 FROM orders WHERE id = ? AND user_id = ? LIMIT 1", orderID, userID).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *orderRepo) CheckShipperOwnership(shipperID, orderID int64) (bool, error) {
	var exists int
	err := r.db.QueryRow("SELECT 1 FROM orders WHERE id = ? AND shipper_id = ? LIMIT 1", orderID, shipperID).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *orderRepo) CountActiveOrdersByShipper(shipperID int64) (int, error) {
	var count int
	// Đếm số đơn đang ship
	err := r.db.QueryRow("SELECT COUNT(*) FROM orders WHERE shipper_id = ? AND order_status = 'shipping'", shipperID).Scan(&count)
	return count, err
}

// 5. LIST & FILTER (HÀM QUAN TRỌNG NHẤT)

func (r *orderRepo) GetOrders(filter map[string]interface{}, page, limit int) ([]models.Order, int, error) {
	offset := (page - 1) * limit

	// Xây dựng câu Query động
	whereClauses := []string{"1=1"}
	var args []interface{}

	if val, ok := filter["user_id"]; ok {
		whereClauses = append(whereClauses, "user_id = ?")
		args = append(args, val)
	}
	if val, ok := filter["shipper_id"]; ok {
		whereClauses = append(whereClauses, "shipper_id = ?")
		args = append(args, val)
	}
	if val, ok := filter["status"]; ok {
		whereClauses = append(whereClauses, "order_status = ?")
		args = append(args, val)
	}
	if val, ok := filter["status_not"]; ok {
		whereClauses = append(whereClauses, "order_status != ?")
		args = append(args, val)
	}

	whereStr := " WHERE " + strings.Join(whereClauses, " AND ")

	// 1. Query Tổng số (Total)
	var total int
	countQuery := "SELECT COUNT(*) FROM orders" + whereStr
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 2. Query Dữ liệu (Data)
	dataQuery := `SELECT id, user_id, shipper_id, payment_status, order_status, latitude, longitude, total_amount, thumbnail_id, created_at, updated_at 
				  FROM orders` + whereStr + ` ORDER BY id DESC LIMIT ? OFFSET ?`

	// Thêm limit/offset vào args
	args = append(args, limit, offset)

	rows, err := r.db.Query(dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		var shipperID sql.NullInt64
		err := rows.Scan(&o.ID, &o.UserID, &shipperID, &o.PaymentStatus, &o.OrderStatus, &o.Latitude, &o.Longitude, &o.TotalAmount, &o.ThumbnailID, &o.CreatedAt, &o.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		if shipperID.Valid {
			o.ShipperID = shipperID.Int64
		}
		orders = append(orders, o)
	}
	return orders, total, nil
}

// 6. STATS & INFO

func (r *orderRepo) GetShipperInfoByOrderID(orderID int64) (int64, string, string, error) {
	var id int64
	var name, phone string
	query := `SELECT u.id, u.name, u.phone 
			  FROM users u 
			  JOIN orders o ON u.id = o.shipper_id 
			  WHERE o.id = ?`
	err := r.db.QueryRow(query, orderID).Scan(&id, &name, &phone)
	if err == sql.ErrNoRows {
		return 0, "", "", nil // Không có shipper cũng không sao
	}
	return id, name, phone, err
}

func (r *orderRepo) GetStats() (int64, float64, error) {
	// Query 1: Đếm đơn (trừ cancelled)
	var numOrder int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM orders WHERE order_status != 'cancelled'").Scan(&numOrder)
	if err != nil {
		return 0, 0, err
	}

	// Query 2: Tổng doanh thu (paid)
	var revenue sql.NullFloat64
	err = r.db.QueryRow("SELECT SUM(total_amount) FROM orders WHERE payment_status = 'paid'").Scan(&revenue)
	if err != nil {
		return 0, 0, err
	}

	return numOrder, revenue.Float64, nil
}

func (r *orderRepo) GetImageURLsByIDs(imageIDs []int64) (map[int64]string, error) {
	if len(imageIDs) == 0 {
		return map[int64]string{}, nil
	}
	// Tạo placeholders (?, ?, ?)
	placeholders := make([]string, len(imageIDs))
	args := make([]interface{}, len(imageIDs))
	for i, id := range imageIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf("SELECT id, url FROM Images WHERE id IN (%s)", strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64]string)
	for rows.Next() {
		var id int64
		var url string
		if err := rows.Scan(&id, &url); err != nil {
			return nil, err
		}
		result[id] = url
	}
	return result, nil
}
