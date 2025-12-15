package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"example.com/delivery-app/dto"
	"example.com/delivery-app/service"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// === CREATE ORDER ===
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID := c.GetInt64("userID")
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.orderService.CreateOrder(userID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully"})
}

// === GET DETAIL ===
func (h *OrderHandler) GetDetail(c *gin.Context) {
	orderID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	userID := c.GetInt64("userID")
	role := c.GetString("role")

	resp, err := h.orderService.GetOrderDetail(orderID, userID, role)
	if err != nil {
		if errors.Is(err, service.ErrNotYourOrder) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else if errors.Is(err, service.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, resp)
}

// === LIST: CUSTOMER MY ORDERS ===
func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	userID := c.GetInt64("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filter := map[string]interface{}{"user_id": userID}

	resp, err := h.orderService.GetListOrders(filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// === LIST: ADMIN ALL ORDERS ===
func (h *OrderHandler) GetAllOrdersAdmin(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	resp, err := h.orderService.GetListOrders(nil, page, limit) // filter rỗng = lấy hết
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// === LIST: SHIPPER AVAILABLE (Processing) ===
func (h *OrderHandler) GetAvailableOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filter := map[string]interface{}{"status": "processing"}

	resp, err := h.orderService.GetListOrders(filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// === LIST: SHIPPER MY ORDERS (Shipping) ===
func (h *OrderHandler) GetMyShippingOrders(c *gin.Context) {
	shipperID := c.GetInt64("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filter := map[string]interface{}{
		"shipper_id": shipperID,
		"status":     "shipping",
	}

	resp, err := h.orderService.GetListOrders(filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// === ACTION: SHIPPER RECEIVE ===
func (h *OrderHandler) ReceiveOrder(c *gin.Context) {
	shipperID := c.GetInt64("userID")
	var req dto.ReceiveOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.orderService.ReceiveOrder(shipperID, req.OrderID)
	if err != nil {
		if errors.Is(err, service.ErrMaxOrdersReached) || errors.Is(err, service.ErrOrderNotProcessing) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order received successfully"})
}

// === ACTION: SHIPPER UPDATE STATUS ===
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	shipperID := c.GetInt64("userID")
	var req dto.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.OrderID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID required"})
		return
	}

	err := h.orderService.UpdateOrder(shipperID, &req)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotOwned) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order updated successfully"})
}

// === ACTION: ADMIN ACCEPT ===
func (h *OrderHandler) AdminAcceptOrder(c *gin.Context) {
	orderID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.orderService.AdminAcceptOrder(orderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order accepted (processing)"})
}

// === ACTION: CANCEL ORDER ===
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID := c.GetInt64("userID")
	orderID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := h.orderService.CancelOrder(userID, orderID); err != nil {
		if errors.Is(err, service.ErrCannotCancel) || errors.Is(err, service.ErrNotYourOrder) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled"})
}

// === STATS ===
func (h *OrderHandler) GetStats(c *gin.Context) {
	count, revenue, err := h.orderService.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"number_of_order": count, "revenue": revenue})
}
