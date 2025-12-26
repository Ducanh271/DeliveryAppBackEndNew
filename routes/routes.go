package routes

import (
	"example.com/delivery-app/handlers"
	"example.com/delivery-app/middleware"
	"example.com/delivery-app/websocket"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	authHandler *handlers.AuthHandler,
	authMiddleware gin.HandlerFunc,
	userHandler *handlers.UserHandler,
	productHandler *handlers.ProductHandler,
	orderHandler *handlers.OrderHandler,
	reviewHandler *handlers.ReviewHandler,
	wsHandler *websocket.WSHandler,
	chatHandler *handlers.ChatHandler, // <--- THÊM MỚI
	rateLimitMiddleware gin.HandlerFunc,
) {
	api := r.Group("api/v1")
	// chống bruce force dò mật khẩu otp
	authGroup := api.Group("/")
	authGroup.Use(rateLimitMiddleware)
	{
		authGroup.POST("/signup", authHandler.SignUp)
		authGroup.POST("/verify-otp", authHandler.VerifyOTPHandler)
		authGroup.POST("/login", authHandler.LoginHandler)
		authGroup.POST("/forgot-password", authHandler.ForgotPasswordHandler)
		authGroup.POST("/verify-reset-otp", authHandler.VerifyResetOTPHandler)
		authGroup.POST("reset-password", authHandler.ResetPasswordHandler)
		authGroup.POST("/refresh-access-token", authHandler.RefreshToken)
		authGroup.POST("/logout", authHandler.Logout)
	}
	protected := api.Group("/")
	protected.Use(authMiddleware)
	protected.GET("/profile", userHandler.Profile)
	// Customer routes
	customer := protected.Group("/customer")
	customer.Use(middleware.RoleMiddleWare("customer")) // Bảo vệ
	customer.POST("/create-order", orderHandler.CreateOrder)
	customer.GET("/orders", orderHandler.GetMyOrders)
	customer.GET("/orders/:id", orderHandler.GetDetail)
	customer.DELETE("/orders/:id", orderHandler.CancelOrder)
	customer.POST("/create-review", reviewHandler.CreateReview)

	// --- Shipper ---
	shipper := protected.Group("/shipper")
	shipper.Use(middleware.RoleMiddleWare("shipper"))

	shipper.POST("/receive-order", orderHandler.ReceiveOrder)
	shipper.POST("/update-order", orderHandler.UpdateOrder)
	shipper.GET("/orders", orderHandler.GetAvailableOrders)                  // Lấy đơn "processing" để nhận
	shipper.GET("/orders/received-orders", orderHandler.GetMyShippingOrders) // Lấy đơn mình đang ship
	shipper.GET("/orders/:id", orderHandler.GetDetail)
	protected.GET("/orders/:id/messages", chatHandler.GetMessages)

	// --- Admin ---
	admin := protected.Group("/admin")
	admin.Use(middleware.RoleMiddleWare("admin")) // Bảo vệ

	admin.POST("/shippers", userHandler.CreateShipper)
	admin.GET("/shippers", userHandler.GetAllShippers)
	admin.POST("/shippers/ban/:id", userHandler.BanUser)
	admin.POST("/shippers/unban/:id", userHandler.UnbanUser)

	admin.GET("/customers", userHandler.GetAllCustomers)
	admin.POST("/customers/ban/:id", userHandler.BanUser)
	admin.POST("/customers/unban/:id", userHandler.UnbanUser)

	// Admin (Create, Delete)
	admin.POST("/products/create-product", productHandler.CreateProduct)
	admin.DELETE("/products/:id", productHandler.DeleteProduct)

	admin.GET("/orders", orderHandler.GetAllOrdersAdmin)
	admin.GET("/orders/:id", orderHandler.GetDetail)
	admin.POST("/orders/accept-order/:id", orderHandler.AdminAcceptOrder)
	admin.GET("/stats", userHandler.GetDashboardStats)
	admin.GET("/stats/orders", orderHandler.GetStats)
	// Public (Get, Search)
	// Lưu ý: Giờ đây Search và Get All dùng chung 1 hàm handler
	api.GET("/products", productHandler.GetProducts)        // Xử lý cả /products?page=1 và /products?q=abc
	api.GET("/products/search", productHandler.GetProducts) // (Optional) Để tương thích ngược nếu FE đang gọi route này
	api.GET("/products/:id", productHandler.GetProductByID)
	api.GET("/products/:id/reviews", reviewHandler.GetReviewsByProduct)
	api.GET("/ws", wsHandler.ServeWs)

}
