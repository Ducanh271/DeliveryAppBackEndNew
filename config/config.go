package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// EmailConfig vẫn giữ nguyên
type EmailConfig struct {
	From     string
	Password string
	Host     string
	Port     string
}

// Thêm DBConfig để gom nhóm
type DBConfig struct {
	User string
	Pass string
	Host string
	Port string
	Name string
	URL  string // Chúng ta sẽ tự tạo DSN (URL)
}

// Config là struct "cha" chứa TẤT CẢ
type Config struct {
	Email         EmailConfig
	DB            DBConfig
	CloudinaryURL string
	JWTSecret     string
	Port          string // Thêm Port cho server
}

// ❗️ Không còn biến toàn cục ở đây

// LoadConfig giờ trả về (Config, error)
func LoadConfig() (Config, error) {
	// Load .env file (Giả sử .env ở thư mục gốc,
	// đường dẫn "../.env" chỉ đúng khi chạy từ package con)
	err := godotenv.Load("../.env") // <-- Sửa đường dẫn nếu cần
	if err != nil {
		// Không fatal, chỉ cảnh báo, để cho phép dùng biến môi trường
		fmt.Println("⚠️  Không tìm thấy file .env, dùng biến môi trường có sẵn")
	}

	// Tạo một biến `cfg` cục bộ
	var cfg Config

	// 1. Nạp Config Email
	cfg.Email = EmailConfig{
		From:     os.Getenv("EMAIL_FROM"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
	}

	// 2. Nạp Config DB
	cfg.DB.User = os.Getenv("DB_USER")
	cfg.DB.Pass = os.Getenv("DB_PASS")
	cfg.DB.Host = os.Getenv("DB_HOST")
	cfg.DB.Port = os.Getenv("DB_PORT")
	cfg.DB.Name = os.Getenv("DB_NAME")
	cfg.DB.URL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DB.User,
		cfg.DB.Pass,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)
	// 3. Nạp các biến còn lại
	cfg.CloudinaryURL = os.Getenv("CLOUDINARY_URL")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	cfg.Port = os.Getenv("PORT") // Dùng PORT cho server

	// --- Kiểm tra và xử lý ---

	// Đặt giá trị mặc định cho Port
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	// Tạo DB Connection String (DSN)
	// (Thay đổi "mysql" nếu bạn dùng driver khác như "postgres")
	cfg.DB.URL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name,
	)

	// Kiểm tra các biến bắt buộc
	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("❌ JWT_SECRET chưa được set trong .env")
	}
	if cfg.CloudinaryURL == "" {
		return Config{}, fmt.Errorf("❌ CLOUDINARY_URL chưa được set trong .env")
	}

	// Trả về struct config đã nạp đầy đủ
	return cfg, nil
}

//
// import (
// 	"fmt"
// 	"log"
// 	"os"
//
// 	"github.com/joho/godotenv"
// )
//
// type EmailConfig struct {
// 	From     string
// 	Password string
// 	Host     string
// 	Port     string
// }
//
// var (
// 	CloudinaryURL string
// 	Email         EmailConfig
// 	JWTSecret     string
// 	DBUser        string
// 	DBPass        string
// 	DBHost        string
// 	DBPort        string
// 	DBName        string
// )
//
// func LoadConfig() {
// 	// Load .env file
// 	err := godotenv.Load("../.env")
// 	if err != nil {
// 		log.Println("⚠️  Không tìm thấy file .env, dùng biến môi trường có sẵn")
// 	}
// 	Email = EmailConfig{
// 		From:     os.Getenv("EMAIL_FROM"),
// 		Password: os.Getenv("EMAIL_PASSWORD"),
// 		Host:     os.Getenv("SMTP_HOST"),
// 		Port:     os.Getenv("SMTP_PORT"),
// 	}
// 	//cloudinary url
// 	CloudinaryURL = os.Getenv("CLOUDINARY_URL")
// 	fmt.Println(CloudinaryURL)
// 	// Gán giá trị
// 	JWTSecret = os.Getenv("JWT_SECRET")
// 	DBUser = os.Getenv("DB_USER")
// 	DBPass = os.Getenv("DB_PASS")
// 	DBHost = os.Getenv("DB_HOST")
// 	DBPort = os.Getenv("DB_PORT")
// 	DBName = os.Getenv("DB_NAME")
//
// 	if JWTSecret == "" {
// 		log.Fatal("❌ JWT_SECRET chưa được set trong .env")
// 	}
// }
