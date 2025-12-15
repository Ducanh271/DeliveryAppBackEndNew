package repository

import (
	"database/sql"
	"example.com/delivery-app/models"
	"fmt"
	"strings"
	// "time"
)

// ProductResponse là DTO "fat query".
// Service sẽ sử dụng các hàm repo "clean" để xây dựng DTO này.
//
//	type ProductResponse struct {
//		ID          int64                 `json:"id"`
//		Name        string                `json:"name"`
//		Description string                `json:"description"`
//		Price       float64               `json:"price"`
//		QtyInitial  int64                 `json:"qty_initial"`
//		QtySold     int64                 `json:"qty_sold"`
//		CreatedAt   time.Time             `json:"created_at"`
//		Images      []models.ProductImage `json:"images"`
//		AvgRate     float64               `json:"avg_rate"`
//		ReviewCount int64                 `json:"review_count"`
//	}
//
// 1. ĐỊNH NGHĨA INTERFACE (Đã sửa)
type ProductRepository interface {
	// Write (Dùng trong UoW)
	Create(product *models.Product) (int64, error)
	CreateImage(imageURL string, publicID string) (int64, error)          // ⬅️ SỬA
	LinkImageToProduct(productID int64, imageID int64, isMain bool) error // ⬅️ SỬA
	DeleteProduct(productID int64) error                                  // ⬅️ SỬA
	DeleteImages(imageIDs []int64) error                                  // ⬅️ SỬA (int -> int64)

	// Read (Không cần Tx)
	GetTotalCount() (int64, error)
	GetByID(id int64) (*models.Product, error)
	GetImagesByProductID(id int64) ([]models.ProductImage, error)
	GetPublicIDsByProductID(productID int64) ([]string, error) // ⬅️ THÊM (Cần cho Service)
	GetAverageRating(id int64) (float64, int64, error)

	// Batch queries
	GetImagesByProductIDs(ids []int64) (map[int64][]models.ProductImage, error)
	GetAverageRatingsByProductIDs(ids []int64) (map[int64]struct {
		Avg   float64
		Count int64
	}, error)

	// Read (Paginated) - Chỉ trả về model
	GetAll(page, limit int64) ([]models.Product, int64, error)               // ⬅️ SỬA (int -> int64)
	Search(query string, page, limit int64) ([]models.Product, int64, error) // ⬅️ SỬA (int -> int64)

	WithTx(tx *sql.Tx) ProductRepository
}

// 2. STRUCT TRIỂN KHAI
type productRepo struct {
	db DBTX
}

// 3. HÀM KHỞI TẠO
func NewProductRepository(db DBTX) ProductRepository {
	return &productRepo{db: db}
}

// 4. TRIỂN KHAI WithTx
func (r *productRepo) WithTx(tx *sql.Tx) ProductRepository {
	return &productRepo{db: tx}
}

// 5. TRIỂN KHAI CÁC METHOD

// === Write Methods ===

func (r *productRepo) Create(product *models.Product) (int64, error) {
	query := "INSERT INTO Products (name, description, price, qty_initial, qty_sold, created_at) VALUES (?, ?, ?, ?, ?, ?)"
	res, err := r.db.Exec(query, product.Name, product.Description, product.Price, product.QtyInitial, product.QtySold, product.CreatedAt)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

// SỬA: Tách từ logic AddProductImageTx cũ
func (r *productRepo) CreateImage(imageURL string, publicID string) (int64, error) {
	queryImg := `INSERT INTO Images (url, public_id) VALUES (?, ?)`
	result, err := r.db.Exec(queryImg, imageURL, publicID)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// SỬA: Tách từ logic AddProductImageTx cũ
func (r *productRepo) LinkImageToProduct(productID int64, imageID int64, isMain bool) error {
	queryMap := `INSERT INTO ProductImages (product_id, image_id, is_main) VALUES (?, ?, ?)`
	_, err := r.db.Exec(queryMap, productID, imageID, isMain)
	return err
}

// SỬA: Thêm check RowsAffected
func (r *productRepo) DeleteProduct(productID int64) error {
	query := "DELETE FROM Products WHERE id = ?"
	result, err := r.db.Exec(query, productID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // Quan trọng: Báo cho service biết
	}
	return nil
}

// SỬA: int -> int64
func (r *productRepo) DeleteImages(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders := "?" + strings.Repeat(",?", len(ids)-1)
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	query := fmt.Sprintf("DELETE FROM Images WHERE id IN (%s)", placeholders)
	_, err := r.db.Exec(query, args...)
	return err
}

// === Read Methods ===

func (r *productRepo) GetTotalCount() (int64, error) {
	var total int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM Products").Scan(&total)
	return total, err
}

func (r *productRepo) GetByID(id int64) (*models.Product, error) {
	query := `SELECT id, name, description, price, qty_initial, qty_sold, created_at
			  FROM Products WHERE id = ?`
	row := r.db.QueryRow(query, id)
	var p models.Product
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.QtyInitial, &p.QtySold, &p.CreatedAt)
	return &p, err
}

// SỬA: Sửa lỗi SQL (bỏ dấu phẩy) và thêm `i.public_id`
func (r *productRepo) GetImagesByProductID(id int64) ([]models.ProductImage, error) {
	query := `
		SELECT i.id, i.url, pi.is_main, i.public_id
		FROM Images i
		INNER JOIN ProductImages pi ON i.id = pi.image_id
		WHERE pi.product_id = ?
	`
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.ProductImage
	for rows.Next() {
		var img models.ProductImage
		var publicID sql.NullString // public_id có thể NULL

		// SỬA: Scan thêm publicID
		err := rows.Scan(&img.ID, &img.URL, &img.IsMain, &publicID)
		if err != nil {
			return nil, err
		}
		if publicID.Valid {
			img.PublicID = publicID.String
		}
		img.ProductID = id // Thêm ProductID vào struct
		images = append(images, img)
	}
	return images, nil
}

// THÊM: Hàm này bị thiếu trong implementation
func (r *productRepo) GetPublicIDsByProductID(productID int64) ([]string, error) {
	rows, err := r.db.Query(`
		SELECT i.public_id 
		FROM Images i 
		JOIN ProductImages pi ON i.id = pi.image_id 
		WHERE pi.product_id = ?`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var publicIDs []string
	for rows.Next() {
		var publicID sql.NullString
		if err := rows.Scan(&publicID); err != nil {
			return nil, err
		}
		if publicID.Valid && publicID.String != "" {
			publicIDs = append(publicIDs, publicID.String)
		}
	}
	return publicIDs, rows.Err()
}

func (r *productRepo) GetAverageRating(id int64) (float64, int64, error) {
	var avgRate sql.NullFloat64
	var count int64
	query := `SELECT IFNULL(AVG(rate), 0), COUNT(*) FROM Reviews WHERE product_id = ?`
	err := r.db.QueryRow(query, id).Scan(&avgRate, &count)
	if err != nil && err != sql.ErrNoRows {
		return 0, 0, err
	}
	return avgRate.Float64, count, nil
}

// === Read Paginated Methods ===
func (r *productRepo) GetAll(page, limit int64) ([]models.Product, int64, error) {
	offset := (page - 1) * limit
	var total int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM Products").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, name, description, price, qty_initial, qty_sold, created_at
			  FROM Products
			  ORDER BY id DESC
			  LIMIT ? OFFSET ?`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.QtyInitial, &p.QtySold, &p.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, p)
	}
	return products, total, nil
}

// SỬA: int -> int64
func (r *productRepo) Search(keyword string, page, limit int64) ([]models.Product, int64, error) {
	offset := (page - 1) * limit
	searchTerm := "%" + keyword + "%"

	var total int64
	countQuery := `SELECT COUNT(*) FROM Products WHERE name LIKE ? OR description LIKE ?`
	err := r.db.QueryRow(countQuery, searchTerm, searchTerm).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, name, description, price, qty_initial, qty_sold, created_at
			  FROM Products
			  WHERE name LIKE ? OR description LIKE ?
			  ORDER BY id DESC
			  LIMIT ? OFFSET ?`
	rows, err := r.db.Query(query, searchTerm, searchTerm, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.QtyInitial, &p.QtySold, &p.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, p)
	}
	return products, total, nil
}

// === Batch Methods (Giải quyết N+1) ===

func (r *productRepo) GetImagesByProductIDs(ids []int64) (map[int64][]models.ProductImage, error) {
	if len(ids) == 0 {
		return map[int64][]models.ProductImage{}, nil
	}

	placeholders := "?" + strings.Repeat(",?", len(ids)-1)
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	// SỬA: Lỗi SQL (dấu phẩy) và lỗi Scan (thiếu cột)
	query := fmt.Sprintf(`
		SELECT pi.product_id, i.id, i.url, pi.is_main, i.public_id
		FROM Images i
		JOIN ProductImages pi ON i.id = pi.image_id
		WHERE pi.product_id IN (%s)`, placeholders)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64][]models.ProductImage)
	for rows.Next() {
		var productID int64
		var img models.ProductImage
		var publicID sql.NullString

		// SỬA: Scan 5 cột
		if err := rows.Scan(&productID, &img.ID, &img.URL, &img.IsMain, &publicID); err != nil {
			return nil, err
		}
		if publicID.Valid {
			img.PublicID = publicID.String
		}
		img.ProductID = productID
		result[productID] = append(result[productID], img)
	}
	return result, rows.Err()
}

func (r *productRepo) GetAverageRatingsByProductIDs(ids []int64) (map[int64]struct {
	Avg   float64
	Count int64
}, error) {
	if len(ids) == 0 {
		return make(map[int64]struct {
			Avg   float64
			Count int64
		}), nil
	}

	placeholders := "?" + strings.Repeat(",?", len(ids)-1)
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT product_id, IFNULL(AVG(rate), 0), COUNT(*)
		FROM Reviews 
		WHERE product_id IN (%s)
		GROUP BY product_id`, placeholders)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64]struct {
		Avg   float64
		Count int64
	})
	for rows.Next() {
		var productID int64
		var avg float64
		var count int64
		if err := rows.Scan(&productID, &avg, &count); err != nil {
			return nil, err
		}
		result[productID] = struct {
			Avg   float64
			Count int64
		}{Avg: avg, Count: count}
	}
	return result, rows.Err()
}
