package database // Khai báo package "database", nơi chứa logic kết nối và quản lý cơ sở dữ liệu

import (
	"log"
	"url-shortener/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Khai báo biến toàn cục DB để có thể sử dụng ở các package khác
var DB *gorm.DB

// Init là hàm khởi tạo kết nối tới cơ sở dữ liệu và thực hiện migrate
func Init() {
	var err error // Khai báo biến err để lưu lỗi nếu có

	dsn := "host=localhost user=postgres password=abc12345 dbname=urlshortener port=5432 sslmode=disable"

	// Mở kết nối tới PostgreSQL thông qua GORM
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// Nếu lỗi xảy ra khi kết nối DB, dừng chương trình và in thông báo lỗi
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Tự động tạo bảng URL trong DB (nếu chưa có) dựa vào struct models.URL
	if err := DB.AutoMigrate(&models.URL{}); err != nil {
		// Nếu có lỗi trong quá trình migrate, dừng chương trình và in thông báo lỗi
		log.Fatalf("Migration failed: %v", err)
	}
}
