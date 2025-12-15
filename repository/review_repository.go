package repository

import (
	"database/sql"
	"example.com/delivery-app/models"
	"fmt"
	"strings"
	"time"
)

type ReviewRepository interface {
	// Write
	Create(review *models.Review) (int64, error)
	CreateImage(imageURL, publicID string) (int64, error)
	LinkImageToReview(reviewID, imageID int64) error

	// Check Logic
	CanUserReview(userID, productID int64) (bool, error)
	CheckReviewExists(userID, orderID, productID int64) (bool, error) // Chặn spam review nhiều lần

	// Read
	GetByProductID(productID int64, page, limit int) ([]models.ReviewDetail, int, error)

	// Batch Query (Giải quyết N+1)
	GetImagesByReviewIDs(reviewIDs []int64) (map[int64][]models.ReviewImage, error)

	WithTx(tx *sql.Tx) ReviewRepository
}

type reviewRepo struct {
	db DBTX
}

func NewReviewRepository(db DBTX) ReviewRepository {
	return &reviewRepo{db: db}
}

func (r *reviewRepo) WithTx(tx *sql.Tx) ReviewRepository {
	return &reviewRepo{db: tx}
}

// --- IMPLEMENTATION ---

// 1. Check Logic
func (r *reviewRepo) CanUserReview(userID, productID int64) (bool, error) {
	// Logic cũ của bạn: User phải có đơn hàng 'paid' & 'delivered' chứa sản phẩm này
	query := `
        SELECT COUNT(*) 
        FROM orders o
        JOIN order_items oi ON o.id = oi.order_id
        WHERE o.user_id = ? 
          AND oi.product_id = ?
          AND o.payment_status = 'paid'
          AND o.order_status = 'delivered'
    `
	var count int
	err := r.db.QueryRow(query, userID, productID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *reviewRepo) CheckReviewExists(userID, orderID, productID int64) (bool, error) {
	query := "SELECT COUNT(*) FROM Reviews WHERE user_id = ? AND order_id = ? AND product_id = ?"
	var count int
	err := r.db.QueryRow(query, userID, orderID, productID).Scan(&count)
	return count > 0, err
}

// 2. Write
func (r *reviewRepo) Create(review *models.Review) (int64, error) {
	query := `INSERT INTO Reviews (product_id, user_id, order_id, rate, content, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	review.CreatedAt = time.Now()
	res, err := r.db.Exec(query, review.ProductID, review.UserID, review.OrderID, review.Rate, review.Content, review.CreatedAt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *reviewRepo) CreateImage(imageURL, publicID string) (int64, error) {
	// Insert vào bảng chung Images
	query := `INSERT INTO Images (url, public_id) VALUES (?, ?)`
	res, err := r.db.Exec(query, imageURL, publicID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *reviewRepo) LinkImageToReview(reviewID, imageID int64) error {
	// Insert vào bảng nối ReviewImages
	query := `INSERT INTO ReviewImages (review_id, image_id) VALUES (?, ?)`
	_, err := r.db.Exec(query, reviewID, imageID)
	return err
}

// 3. Read List
func (r *reviewRepo) GetByProductID(productID int64, page, limit int) ([]models.ReviewDetail, int, error) {
	offset := (page - 1) * limit

	// Query 1: Total
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM Reviews WHERE product_id = ?", productID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Query 2: Data (Join Users để lấy tên)
	query := `
		SELECT r.id, r.product_id, r.user_id, r.order_id, r.rate, r.content, r.created_at, u.name
		FROM Reviews r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.product_id = ?
		ORDER BY r.created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.Query(query, productID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reviews []models.ReviewDetail
	for rows.Next() {
		var rd models.ReviewDetail
		err := rows.Scan(
			&rd.ID, &rd.ProductID, &rd.UserID, &rd.OrderID, &rd.Rate, &rd.Content, &rd.CreatedAt, &rd.UserName,
		)
		if err != nil {
			return nil, 0, err
		}
		reviews = append(reviews, rd)
	}
	return reviews, total, nil
}

// 4. Batch Query Images
func (r *reviewRepo) GetImagesByReviewIDs(reviewIDs []int64) (map[int64][]models.ReviewImage, error) {
	if len(reviewIDs) == 0 {
		return map[int64][]models.ReviewImage{}, nil
	}

	placeholders := "?" + strings.Repeat(",?", len(reviewIDs)-1)
	args := make([]interface{}, len(reviewIDs))
	for i, id := range reviewIDs {
		args[i] = id
	}

	// Join 2 bảng: ReviewImages và Images
	query := fmt.Sprintf(`
		SELECT ri.review_id, i.id, i.url, i.public_id
		FROM Images i
		JOIN ReviewImages ri ON i.id = ri.image_id
		WHERE ri.review_id IN (%s)
	`, placeholders)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64][]models.ReviewImage)
	for rows.Next() {
		var reviewID int64
		var img models.ReviewImage
		var publicID sql.NullString

		if err := rows.Scan(&reviewID, &img.ID, &img.URL, &publicID); err != nil {
			return nil, err
		}
		if publicID.Valid {
			img.PublicID = publicID.String
		}
		img.ReviewID = reviewID

		result[reviewID] = append(result[reviewID], img)
	}
	return result, nil
}
