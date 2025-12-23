package main

import (
	"example.com/delivery-app/config"
	"example.com/delivery-app/handlers"
	"example.com/delivery-app/infrastructure/storage"
	"example.com/delivery-app/middleware"
	"example.com/delivery-app/notification"
	"example.com/delivery-app/repository"
	"example.com/delivery-app/service"
	"example.com/delivery-app/websocket"
	"github.com/cloudinary/cloudinary-go/v2"

	// "github.com/cloudinary/cloudinary-go/v2"
	"log"
	"net/http"

	"example.com/delivery-app/database"
	"example.com/delivery-app/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}
	// Init DB
	db, err := database.NewConnection(cfg.DB.URL)
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Close()
	if err := database.CreateDefaultAdmin(db); err != nil {
		log.Fatal("Error seeding admin:", err)
	}
	cld, err := cloudinary.NewFromURL(cfg.CloudinaryURL)
	if err != nil {
		log.Fatal("Failed to connect to Cloudinary")
	}
	// 4. === KHỞI TẠO (WIRING) CÁC TẦNG ===
	// (Đây là phần quan trọng nhất của DI)

	// --- Tầng Infrastructure (Các dịch vụ bên ngoài) ---
	// (Giả sử cfg.Email là struct config cho email)
	emailService := notification.NewEmailService(cfg.Email)
	imgSvc := storage.NewCloudinaryService(cld)
	// --- Tầng Repository (và Unit of Work) ---
	unitOfWork := repository.NewUnitOfWork(db)   // UoW (cho Write)
	userRepo := repository.NewUserRepository(db) // Repo gốc (cho Read)
	productRepo := repository.NewProductRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	msgRepo := repository.NewMessageRepository(db)

	chatService := service.NewChatService(msgRepo, orderRepo)

	hub := websocket.NewHub(chatService)
	go hub.Run() // Chạy Hub ở background

	wsHandler := websocket.NewWSHandler(hub, cfg.JWTSecret)
	// Khoi chay job don dep
	repository.StartTokenCleanUp(db)

	// --- Tầng Service (Logic nghiệp vụ) ---
	// Tiêm repo, uow, và email service vào Auth Service
	authService := service.NewAuthService(
		userRepo,
		refreshTokenRepo,
		unitOfWork,
		emailService,
		cfg.JWTSecret,
	)
	// --- Khởi tạo Service MỚI ---
	userService := service.NewUserService(
		userRepo,
		refreshTokenRepo, // Cần cho BanUser
		unitOfWork,       // Cần cho CreateShipper
	)
	productService := service.NewProductService(
		productRepo,
		imgSvc,
		unitOfWork,
	)
	orderService := service.NewOrderService(
		orderRepo,
		productRepo,
		userRepo,
		unitOfWork,
	)
	reviewService := service.NewReviewService(
		reviewRepo,
		imgSvc,
		unitOfWork,
	)
	// --- Khởi tạo Handler MỚI ---
	userHandler := handlers.NewUserHandler(userService) // ⬅️ DÙNG Ở ĐÂY
	// --- Tầng Handler (Giao tiếp với Gin) ---
	// Tiêm auth service vào auth handler
	authHandler := handlers.NewAuthHandler(authService)
	// productHandler := handlers.NewProductHandler(productService) // (Tương lai)
	productHandler := handlers.NewProductHandler(productService)
	// ... các handler khác ...
	orderHandler := handlers.NewOrderHandler(orderService)

	reviewHandler := handlers.NewReviewHandler(reviewService)

	chatHandler := handlers.NewChatHandler(chatService)
	// 4.5 khoi tao middle ware
	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)
	// 5. Khởi tạo Gin Engine
	r := gin.Default()

	// 6. Cấu hình Middleware (CORS)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// 7. Setup Routes (Bản "Clean")
	// Chỉ truyền các HANDLER đã được khởi tạo
	// Không còn truyền `db` hay `cld` vào routes nữa!
	routes.SetupRoutes(r,
		authHandler,
		authMiddleware,
		userHandler,
		productHandler,
		orderHandler,
		reviewHandler,
		wsHandler,
		chatHandler,
	)

	// 8. Run Server
	log.Println("Server is running on port :8080")
	r.Run(":8080")
}

//
// func main() {
// 	// Init DB
// 	config.LoadConfig()
// 	database.InitDB()
// 	defer database.DB.Close()
// 	if err := database.CreateDefaultAdmin(database.DB); err != nil {
// 		log.Fatal("Error seeding admin:", err)
// 	}
//
// 	// Tạo Gin engine
// 	r := gin.Default()
// 	// create cloudinary
// 	cld, err := cloudinary.NewFromURL(config.CloudinaryURL)
// 	if err!= nil {
// 		log.Fatal("Failed to connect to Cloudinary")
// 	}
// 	// Middleware CORS
// 	r.Use(func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
// 		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
//
// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(http.StatusNoContent)
// 			return
// 		}
//
// 		c.Next()
// 	})
//
// 	// Setup routes (truyền DB vào nếu cần)
// 	routes.SetupRoutes(r, database.DB, cld)
//
// 	// Run server
// 	r.Run(":8080")
// }
