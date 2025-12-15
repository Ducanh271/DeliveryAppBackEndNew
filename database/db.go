package database

import (
	"database/sql"
	"fmt"
	"log"

	// Vẫn import driver
	_ "github.com/go-sql-driver/mysql"
)

// ❗️ KHÔNG CÒN BIẾN `var DB *sql.DB`

// InitDB() được đổi tên thành NewConnection()
// Nó NHẬN DSN làm tham số
// Nó TRẢ VỀ (*sql.DB, error)
func NewConnection(dsn string) (*sql.DB, error) {
	// ❗️ Không còn tự tạo DSN, vì nó được truyền vào

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Database connection failed:", err)
		return nil, err
	}

	// Kiểm tra kết nối
	if err := db.Ping(); err != nil {
		log.Fatal("❌ Database not reachable:", err)
		return nil, err
	}

	fmt.Println("✅ Connected to MySQL!")

	// Trả về đối tượng CSDL, không gán vào biến toàn cục
	return db, nil
}

//
// import (
// 	"database/sql"
// 	"example.com/delivery-app/config"
// 	"fmt"
// 	"log"
//
// 	_ "github.com/go-sql-driver/mysql"
// )
//
// var DB *sql.DB
//
// func InitDB() {
//
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
// 		config.DBUser,
// 		config.DBPass,
// 		config.DBHost,
// 		config.DBPort,
// 		config.DBName,
// 	)
// 	var err error
// 	DB, err = sql.Open("mysql", dsn)
// 	if err != nil {
// 		log.Fatal("❌ Database connection failed:", err)
// 	}
//
// 	// Kiểm tra kết nối
// 	if err := DB.Ping(); err != nil {
// 		log.Fatal("❌ Database not reachable:", err)
// 	}
//
// 	fmt.Println("✅ Connected to MySQL!")
// }
