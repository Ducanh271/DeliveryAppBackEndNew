package service

import (
	"example.com/delivery-app/models"
	"example.com/delivery-app/repository"
	"log"
)

type ChatService struct {
	msgRepo   repository.MessageRepository
	orderRepo repository.OrderRepository // Cần để check order tồn tại
}

func NewChatService(msgRepo repository.MessageRepository, orderRepo repository.OrderRepository) *ChatService {
	return &ChatService{msgRepo: msgRepo, orderRepo: orderRepo}
}

// HandleChatMessage: Lưu tin nhắn và trả về tin nhắn đã lưu (có ID, CreatedAt)
func (s *ChatService) HandleChatMessage(msg *models.Message) (*models.Message, error) {
	// 1. Validate (VD: check order có tồn tại không, user có thuộc order không)
	// (Tạm bỏ qua để đơn giản, nhưng nên thêm vào)

	// 2. Lưu vào DB
	if err := s.msgRepo.SaveMessage(msg); err != nil {
		log.Printf("Lỗi lưu tin nhắn: %v", err)
		return nil, err
	}
	return msg, nil
}

// GetReceiverID: Tìm người nhận dựa trên OrderID và Người gửi
// (Ví dụ: Nếu người gửi là Customer, người nhận là Shipper của đơn đó)
func (s *ChatService) GetReceiverID(orderID, senderID int64) (int64, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return 0, err
	}

	if senderID == order.UserID {
		return order.ShipperID, nil
	} else if senderID == order.ShipperID {
		return order.UserID, nil
	}
	return 0, nil // Sender không liên quan đơn hàng
}
