package websocket

import (
	"encoding/json"
	"example.com/delivery-app/models"
	"example.com/delivery-app/service"
	"sync"
)

type Hub struct {
	// Clients: map[UserID]*Client
	// Một user có thể có nhiều kết nối (đa thiết bị), nên value nên là map hoặc slice
	// Nhưng để đơn giản như code cũ, ta giữ 1 connection/user
	Clients map[int64]*Client

	Register   chan *Client
	Unregister chan *Client

	// Kênh nhận tin nhắn từ Client gửi lên
	Broadcast chan *MessageWrapper

	chatService *service.ChatService // ⬅️ Tiêm Service
	mu          sync.RWMutex
}

// Wrapper để biết tin nhắn đến từ ai
type MessageWrapper struct {
	Sender  *Client
	Message *models.Message
}

func NewHub(chatService *service.ChatService) *Hub {
	return &Hub{
		Clients:     make(map[int64]*Client),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan *MessageWrapper),
		chatService: chatService,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.ID] = client // Lưu client theo UserID
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if existing, ok := h.Clients[client.ID]; ok && existing == client {
				delete(h.Clients, client.ID)
				close(client.Send)
			}
			h.mu.Unlock()

		case wrapper := <-h.Broadcast:
			// Xử lý tin nhắn nhận được
			h.handleMessage(wrapper.Sender, wrapper.Message)
		}
	}
}

// Logic xử lý tin nhắn (Tách ra cho gọn)
func (h *Hub) handleMessage(sender *Client, msg *models.Message) {
	// --- 1. KIỂM TRA RATE LIMIT (TỐC ĐỘ) ---
	// Hàm Allow() trả về false nếu hết lượt
	if !sender.Limiter.Allow() {
		h.sendError(sender, "Bạn đang gửi tin quá nhanh. Vui lòng chậm lại!")
		return
	}

	// --- 2. KIỂM TRA ĐỘ DÀI TIN NHẮN ---
	const MaxContentLength = 1000 // Ví dụ: giới hạn 1000 ký tự
	if len(msg.Content) > MaxContentLength {
		h.sendError(sender, "Tin nhắn quá dài (tối đa 1000 ký tự).")
		return
	}
	if len(msg.Content) == 0 && msg.Type == "chat_message" {
		h.sendError(sender, "Nội dung tin nhắn không được để trống.")
		return
	}
	// Xác định người nhận
	var receiverID int64
	var err error

	switch msg.Type {
	case "chat_message":
		// Nếu client không gửi ToUserID, ta tự tìm
		receiverID, err = h.chatService.GetReceiverID(msg.OrderID, sender.ID)
		if err != nil {
			// SỬA: Gửi lỗi về cho sender thay vì return im lặng
			h.sendError(sender, "Không thể gửi tin: "+err.Error())
			return
		}

		if msg.ToUserID == sender.ID {
			h.sendError(sender, "Không thể tự chat với chính mình")
			return
		}
		msg.FromUserID = sender.ID
		savedMsg, err := h.chatService.HandleChatMessage(msg)
		if err != nil {
			h.sendError(sender, "Lỗi: "+err.Error())
			return
		}

		// Gửi cho người nhận
		h.SendToUser(receiverID, savedMsg)

		// (Optional) Gửi lại cho người gửi để xác nhận (ack)
		// h.SendToUser(sender.ID, savedMsg)

	case "location_update":
		// Logic tương tự, tìm người nhận và gửi, không cần lưu DB (hoặc lưu Redis)
		if msg.ToUserID == 0 {
			receiverID, err = h.chatService.GetReceiverID(msg.OrderID, sender.ID)
			if err != nil {
				return
			}
		}
		h.SendToUser(receiverID, msg)
	}
}

func (h *Hub) sendError(client *Client, errMsg string) {
	errorMsg := &models.Message{
		Type:    "error",
		Content: errMsg,
		// Có thể thêm CreatedAt để client hiển thị đúng chuẩn
	}
	data, _ := json.Marshal(errorMsg)
	// Gửi thẳng vào kênh Send của client đó
	select {
	case client.Send <- data:
	default:
		close(client.Send)
		delete(h.Clients, client.ID)
	}
}
func (h *Hub) SendToUser(userID int64, msg *models.Message) {
	h.mu.RLock()
	client, ok := h.Clients[userID]
	h.mu.RUnlock()

	if !ok {
		return
	} // User offline

	data, _ := json.Marshal(msg)
	select {
	case client.Send <- data:
	default:
		// Kết nối chết, đóng và xóa
		h.Unregister <- client
	}
}
