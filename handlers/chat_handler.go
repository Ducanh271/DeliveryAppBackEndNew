// handlers/chat_handler.go
package handlers

import (
	"net/http"
	"strconv"

	"example.com/delivery-app/service"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

// GET /api/v1/orders/:id/messages
func (h *ChatHandler) GetMessages(c *gin.Context) {
	// 1. Lấy Order ID từ URL
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// 2. Lấy params phân trang (limit, before_id)
	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)

	beforeIDStr := c.DefaultQuery("before_id", "0")
	beforeID, _ := strconv.ParseInt(beforeIDStr, 10, 64)

	// 3. Gọi Service
	messages, err := h.chatService.GetMessageHistory(orderID, limit, beforeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Trả về kết quả
	c.JSON(http.StatusOK, gin.H{"data": messages})
}
