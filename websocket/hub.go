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
	// Xác định người nhận
	var receiverID int64
	var err error

	switch msg.Type {
	case "chat_message":
		// Nếu client không gửi ToUserID, ta tự tìm
		if msg.ToUserID == 0 {
			receiverID, err = h.chatService.GetReceiverID(msg.OrderID, sender.ID)
			if err != nil || receiverID == 0 {
				return
			}
			msg.ToUserID = receiverID
		} else {
			receiverID = msg.ToUserID
		}

		// Lưu DB qua Service
		msg.FromUserID = sender.ID
		savedMsg, err := h.chatService.HandleChatMessage(msg)
		if err != nil {
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
