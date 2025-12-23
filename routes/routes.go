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
) {
	api := r.Group("api/v1")
	api.POST("/signup", authHandler.SignUp)
	api.POST("/verify-otp", authHandler.VerifyOTPHandler)
	api.POST("/login", authHandler.LoginHandler)
	api.POST("/forgot-password", authHandler.ForgotPasswordHandler)
	api.POST("/verify-reset-otp", authHandler.VerifyResetOTPHandler)
	api.POST("reset-password", authHandler.ResetPasswordHandler)
	api.POST("/refresh-access-token", authHandler.RefreshToken)
	api.POST("/logout", authHandler.Logout)

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

//
// func SetupRoutes(r *gin.Engine, db *sql.DB, cld *cloudinary.Cloudinary) {
// 	var Hub = websocket.NewHub(db)
// 	go Hub.Run()
// 	api := r.Group("/api/v1")
// 	api.GET("/ws", func(c *gin.Context) {
// 		websocket.ServeWs(Hub, c)
// 	})
// 	// User routes
// 	api.POST("/signup", func(c *gin.Context) {
// 		handlers.SignupHandler(c, db)
// 	})
// 	api.POST("/login", func(c *gin.Context) {
// 		handlers.LoginHandler(c, db)
// 	})
// 	api.POST("logout", func(c *gin.Context) {
// 		handlers.LogoutHandler(c, db)
// 	})
// 	api.POST("/refresh-access-token", func(c *gin.Context) {
// 		handlers.RefreshTokenHandler(c, db)
// 	})
// 	api.POST("/verify-otp", func(c *gin.Context) {
// 		handlers.VerifyOTPHandler(c, db)
// 	})
// 	products := api.Group("/products")
// 	{
// 		products.GET("", func(c *gin.Context) {
// 			handlers.GetProductsHandler(c, db)
// 		})
// 		products.GET("/:id", func(c *gin.Context) {
// 			handlers.GetProductByIDHandler(c, db)
// 		})
// 		products.GET("/:id/reviews", func(c *gin.Context) {
// 			handlers.GetReviewsByProductIDHandler(c, db)
// 		})
// 		products.GET("/search", func(c *gin.Context) {
// 			handlers.SearchProductHandler(c, db)
// 		})
// 	}
//
// 	api.POST("/forgot-password", func(c *gin.Context) { handlers.ForgetPasswordHandler(c, db) })
// 	api.POST("/verify-otp-for-reset", func(c *gin.Context) { handlers.VerifyOTPForResetHandler(c, db) })
// 	api.POST("/reset-password", func(c *gin.Context) { handlers.ResetPasswordHandler(c, db) })
// 	// Profile (bảo vệ bằng JWT)
// 	protected := api.Group("/")
// 	protected.Use(middleware.AuthMiddleware())
// 	protected.GET("/profile", middleware.RoleMiddleWare("customer", "shipper"), func(c *gin.Context) {
// 		handlers.ProfileHandler(c, db)
// 	})
// 	// chỉ cho customer
// 	protected.POST("/create-order", middleware.RoleMiddleWare("customer"), func(c *gin.Context) {
// 		handlers.CreateOrderHandler(c, db)
// 	})
// 	protected.GET("/orders", middleware.RoleMiddleWare("customer"), func(c *gin.Context) {
// 		handlers.GetOrdersByUserIDHandler(c, db)
// 	})
// 	protected.GET("/orders/:id", middleware.RoleMiddleWare("customer", "admin", "shipper"), func(c *gin.Context) {
// 		handlers.GetOrderDetailHandler(c, db)
// 	})
// 	protected.GET("/orders/shipper-info/:id", middleware.RoleMiddleWare("customer"), func(c *gin.Context) {
// 		handlers.GetShipperInfoByOrderIDHandler(c, db)
// 	})
// 	protected.POST("/create-review", middleware.RoleMiddleWare("customer"), func(c *gin.Context) {
// 		handlers.CreateNewReviewHandler(c, db)
// 	})
// 	// Lấy tin nhắn theo đơn, hỗ trợ phân trang
// 	protected.GET("/orders/:id/messages", middleware.RoleMiddleWare("customer", "admin", "shipper"), func(c *gin.Context) {
// 		handlers.GetMessageHandler(c, db)
// 	})
// 	protected.DELETE("/orders/:id", middleware.RoleMiddleWare("customer"), func(c *gin.Context) {
// 		handlers.CancleOrderByUserHandler(c, db)
// 	})
//
// 	// chi cho admin
// 	protected.POST("/admin/create-shipper", middleware.RoleMiddleWare("admin"), func(c *gin.Context) {
// 		handlers.CreateShipper(c, db)
// 	})
// 	protected.POST("/admin/products/create-product", middleware.RoleMiddleWare("admin"), func(c *gin.Context) {
// 		handlers.CreateNewProductHandler(c, db)
// 	})
// 	protected.DELETE("/admin/products/:id", middleware.RoleMiddleWare("admin"), func(c *gin.Context) {
// 		handlers.DeleteProductHandler(c, db, cld)
// 	})
// 	protected.POST("/admin/orders/accept-order/:id", middleware.RoleMiddleWare("admin"), func(c *gin.Context) {
// 		handlers.AcceptOrderAdmin(c, db)
// 	})
// 	protected.GET("/admin/customers", middleware.RoleMiddleWare("admin"), func(c *gin.Context) {
// 		handlers.GetAllCustomersHandler(c, db)
// 	})
// 	protected.GET("/admin/shippers", middleware.RoleMiddleWare("admin"), func(c *gin.Context) {
// 		handlers.GetAllShippersHandler(c, db)
// 	})
// 	protected.POST("/admin/shippers/ban-shipper/:id", middleware.RoleMiddleWare("admin"), func(c *gin.Context) {
// 		handlers.BanUserAccountHandler(c, db)
// 	})
//
// 	protected.POST("/admin/shippers/unban-shipper/:id", middleware.RoleMiddleWare("admin"), func(c *gin.Context) {
// 		handlers.UnBanUserAccountHandler(c, db)
// 	})
//
// 	protected.POST("/admin/customers/ban-customer/:id", middleware.RoleMiddleWare("admin"), func(c *gin.Context) {
// 		handlers.BanUserAccountHandler(c, db)
// 	})
// 	protected.POST("/admin/customers/unban-customer/:id", middleware.RoleMiddleWare("admin"), func(c *gin.Context) {
// 		handlers.UnBanUserAccountHandler(c, db)
// 	})
//
// 	protected.GET("/admin/orders", middleware.RoleMiddleWare("admin"), func(c *gin.Context) {
// 		handlers.GetOrdersByAdminHandler(c, db)
// 	})
// 	// chi cho shipper
// 	protected.POST("/shipper/receive-order", middleware.RoleMiddleWare("shipper"), func(c *gin.Context) {
// 		handlers.ReceiveOrderByShipperHandler(c, db, Hub)
// 	})
// 	protected.POST("/shipper/update-order", middleware.RoleMiddleWare("shipper"), func(c *gin.Context) {
// 		handlers.UpdateOrderShipper(c, db)
// 	})
// 	protected.GET("/shipper/orders", middleware.RoleMiddleWare("shipper"), func(c *gin.Context) {
// 		handlers.GetOrdersByShipperHandler(c, db)
// 	})
// 	protected.GET("/shipper/orders/received-orders", middleware.RoleMiddleWare("shipper"), func(c *gin.Context) {
// 		handlers.GetReceivedOrdersByShipperHandler(c, db)
// 	})
// 	protected.GET("/admin/orders/num-revenue", middleware.RoleMiddleWare("admin"), func(c *gin.Context) {
// 		handlers.GetNumberOfOrderAndRevenueHandler(c, db)
// 	})
// 	protected.GET("/admin/customers/num-customer", middleware.RoleMiddleWare("admin"), func(c *gin.Context) {
// 		handlers.GetNumberOfCustomerHandler(c, db)
// 	})
// 	protected.GET("/admin/shippers/num-shippers", middleware.RoleMiddleWare("admin"), func(c *gin.Context) {
// 		handlers.GetNumberOfShipperHandler(c, db)
// 	})
// 	protected.GET("/admin/products/num-products", middleware.RoleMiddleWare("admin"), func(c *gin.Context) {
// 		handlers.GetNumberOfProductHandler(c, db)
// 	})
// }
