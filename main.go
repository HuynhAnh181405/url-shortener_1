package main // Đây là package chính, là entry point của chương trình Go

import (
	"url-shortener/database" // Import package để khởi tạo kết nối cơ sở dữ liệu
	"url-shortener/handlers" // Import các hàm xử lý HTTP request (shorten, redirect, stats)

	"github.com/gin-gonic/gin" // Gin là web framework dùng để xây dựng RESTful API
)

func main() {
	// ✅ Khởi tạo cơ sở dữ liệu (kết nối PostgreSQL, tự động migrate)
	database.Init()

	// ✅ Tạo một instance mặc định của Gin (có sẵn middleware logger, recovery)
	r := gin.Default()

	// ✅ Định nghĩa route POST /shorten
	// Client gửi long_url → server trả về short_url
	r.POST("/shorten", handlers.ShortenURL)

	// ✅ Định nghĩa route GET /:code
	// Người dùng truy cập short_url → server tìm long_url và chuyển hướng
	r.GET("/:code", handlers.RedirectURL)

	// ✅ Định nghĩa route GET /stats/:code
	// Trả về số lượt click và thông tin long_url tương ứng với short_url
	r.GET("/stats/:code", handlers.GetStats)

	// ✅ Chạy server tại cổng 8080
	r.Run(":8080") // Nếu muốn chạy trên cổng khác, đổi thành ":<port>"
}
