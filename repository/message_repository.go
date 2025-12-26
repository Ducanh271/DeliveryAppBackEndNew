package repository

import (
	"database/sql"
	"example.com/delivery-app/models"
	"example.com/delivery-app/security"
	"time"
)

type MessageRepository interface {
	SaveMessage(msg *models.Message) error
	GetMessagesByOrder(orderID int64, limit int, beforeID int64) ([]models.Message, error)
	GetUnreadCount(userID, orderID int64) (int, error)
	MarkAsRead(userID, orderID int64) error
}

type messageRepo struct {
	db        DBTX
	secretKey string
}

func NewMessageRepository(db DBTX, secretKey string) MessageRepository {
	return &messageRepo{
		db:        db,
		secretKey: secretKey,
	}
}

func (r *messageRepo) SaveMessage(msg *models.Message) error {
	// 1. Mã hóa nội dung tin nhắn
	encryptedContent, err := security.Encrypt(msg.Content, r.secretKey)
	if err != nil {
		return err
	}
	query := `INSERT INTO messages(order_id, sender_id, receiver_id, content, is_read, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	msg.CreatedAt = time.Now()
	res, err := r.db.Exec(query, msg.OrderID, msg.FromUserID, msg.ToUserID, encryptedContent, false, msg.CreatedAt)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	msg.ID = id
	return nil
}

func (r *messageRepo) GetMessagesByOrder(orderID int64, limit int, beforeID int64) ([]models.Message, error) {
	var rows *sql.Rows
	var err error

	// (Logic phân trang cursor-based giống cũ)
	if beforeID == 0 {
		rows, err = r.db.Query(`SELECT id, order_id, sender_id, receiver_id, content, is_read, created_at FROM messages WHERE order_id = ? ORDER BY id DESC LIMIT ?`, orderID, limit)
	} else {
		rows, err = r.db.Query(`SELECT id, order_id, sender_id, receiver_id, content, is_read, created_at FROM messages WHERE order_id = ? AND id < ? ORDER BY id DESC LIMIT ?`, orderID, beforeID, limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var m models.Message
		rows.Scan(&m.ID, &m.OrderID, &m.FromUserID, &m.ToUserID, &m.Content, &m.IsRead, &m.CreatedAt)
		decryptedContent, err := security.Decrypt(m.Content, r.secretKey)
		if err == nil {
			m.Content = decryptedContent // Gán lại nội dung đã giải mã
		} else {
			// Nếu lỗi (ví dụ dữ liệu cũ chưa mã hóa), giữ nguyên hoặc báo lỗi
			m.Content = "[Encrypted Message]"
		}
		m.Type = "message"
		messages = append(messages, m)
	}
	return messages, nil
}

func (r *messageRepo) GetUnreadCount(userID, orderID int64) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM messages WHERE receiver_id = ? AND order_id = ? AND is_read = false", userID, orderID).Scan(&count)
	return count, err
}

func (r *messageRepo) MarkAsRead(userID, orderID int64) error {
	_, err := r.db.Exec("UPDATE messages SET is_read = true WHERE receiver_id = ? AND order_id = ?", userID, orderID)
	return err
}
