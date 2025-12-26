package service

import (
	"errors"
	"example.com/delivery-app/models"
	"example.com/delivery-app/repository"
	"example.com/delivery-app/utils" // <--- Import package utils vừa tạo
	"log"
)

type ChatService struct {
	msgRepo   repository.MessageRepository
	orderRepo repository.OrderRepository // Cần để check order tồn tại
}

func NewChatService(msgRepo repository.MessageRepository, orderRepo repository.OrderRepository) *ChatService {
	return &ChatService{msgRepo: msgRepo, orderRepo: orderRepo}
}
func (s *ChatService) GetReceiverID(orderID, senderID int64) (int64, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return 0, err // Lỗi DB hoặc không tìm thấy order
	}
	if order == nil {
		return 0, ErrOrderNotFound
	}

	// Logic xác định người nhận đối diện
	if senderID == order.UserID {
		// Nếu shipper chưa nhận đơn (ShipperID = 0), khách không thể chat
		if order.ShipperID == 0 {
			return 0, errors.New("đơn hàng chưa có tài xế nhận, không thể chat")
		}
		return order.ShipperID, nil
	}

	if senderID == order.ShipperID {
		return order.UserID, nil
	}

	// Nếu sender không phải Khách, cũng không phải Shipper của đơn này
	return 0, ErrOrderNotOwned
}

// HandleChatMessage: Lưu tin nhắn và trả về tin nhắn đã lưu (có ID, CreatedAt)
func (s *ChatService) HandleChatMessage(msg *models.Message) (*models.Message, error) {
	isOwner, _ := s.orderRepo.CheckOrderOwnership(msg.FromUserID, msg.OrderID)
	if !isOwner {
		isOwner, _ = s.orderRepo.CheckShipperOwnership(msg.FromUserID, msg.OrderID)
		if !isOwner {
			return nil, ErrNotYourOrder
		}
	}

	receiverID, err := s.GetReceiverID(msg.OrderID, msg.FromUserID)
	if err != nil {
		return nil, err
	}

	// Đảm bảo msg.ToUserID khớp với logic hệ thống
	if msg.ToUserID != 0 && msg.ToUserID != receiverID {
		return nil, ErrInvalidReceiver
	}
	msg.ToUserID = receiverID

	msg.Content = utils.SanitizeText(msg.Content)

	// Kiểm tra lại sau khi sanitize (trường hợp user gửi toàn thẻ script thì kết quả sẽ là rỗng)
	if msg.Content == "" {
		return nil, errors.New("nội dung tin nhắn không hợp lệ")
	}

	// 2. Lưu vào DB
	if err := s.msgRepo.SaveMessage(msg); err != nil {
		log.Printf("Lỗi lưu tin nhắn: %v", err)
		return nil, err
	}
	return msg, nil
}

// GetReceiverID: Tìm người nhận dựa trên OrderID và Người gửi
// (Ví dụ: Nếu người gửi là Customer, người nhận là Shipper của đơn đó)
//
//	func (s *ChatService) GetReceiverID(orderID, senderID int64) (int64, error) {
//		order, err := s.orderRepo.GetByID(orderID)
//		if err != nil {
//			return 0, err
//		}
//
//		if senderID == order.UserID {
//			return order.ShipperID, nil
//		} else if senderID == order.ShipperID {
//			return order.UserID, nil
//		}
//		return 0, nil // Sender không liên quan đơn hàng
//	}
func (s *ChatService) GetMessageHistory(orderID int64, limit int, beforeID int64) ([]models.Message, error) {
	// Mặc định limit nếu không truyền
	if limit <= 0 {
		limit = 20
	}
	// Gọi xuống Repository
	return s.msgRepo.GetMessagesByOrder(orderID, limit, beforeID)
}
