// handlers/url.go
package handlers // Khai báo package chứa các hàm xử lý (handler) cho các route/API

import (
	"math/rand"
	"net/http"
	"time"
	"url-shortener/database"
	"url-shortener/models"

	"github.com/gin-gonic/gin"
)

// Bộ ký tự dùng để tạo chuỗi ngắn (short URL)
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Hàm generateShortURL tạo chuỗi ngắn ngẫu nhiên có độ dài `length`
func generateShortURL(length int) string {
	rand.Seed(time.Now().UnixNano()) // Seed bằng thời gian hiện tại để đảm bảo kết quả random khác nhau mỗi lần gọi.
	b := make([]byte, length)        // Tạo slice byte có độ dài bằng `length`
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))] // Gán mỗi phần tử là một ký tự ngẫu nhiên từ charset
	}
	return string(b) // Chuyển slice byte sang string
}

// API POST /shorten – Nhận long_url và trả về short_url
func ShortenURL(c *gin.Context) {
	// Định nghĩa struct tạm thời để parse JSON đầu vào
	var request struct {
		LongURL string `json:"long_url"` // Trường JSON là "long_url"
	}

	// Kiểm tra dữ liệu đầu vào có hợp lệ không
	// c.ShouldBindJSON(&request) → đọc dữ liệu JSON từ request body và lưu vào biến request
	// err != nil || request.LongURL == "", gửi lỗi hoặc để trống
	if err := c.ShouldBindJSON(&request); err != nil || request.LongURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"}) // Trả về lỗi 400 nếu dữ liệu sai
		return
	}

	// Kiểm tra nếu long_url đã tồn tại trong DB thì trả lại short_url cũ
	var existing models.URL
	if err := database.DB.Where("long_url = ?", request.LongURL).First(&existing).Error; err == nil { //nếu không có lỗi (tức là tìm thấy), thì err == nil → điều kiện if đúng.
		c.JSON(http.StatusOK, gin.H{"short_url": existing.ShortURL}) //tìm thấy URL trong DB rồi → trả về chuỗi short URL cũ đã lưu trước đó.

		return
	}

	// Sinh short URL ngẫu nhiên và kiểm tra xem có bị trùng không
	var short string
	for {
		short = generateShortURL(6) // Tạo short URL dài 6 ký tự
		var temp models.URL
		// Kiểm tra xem short URL đã tồn tại trong DB chưa
		if err := database.DB.Where("short_url = ?", short).First(&temp).Error; err != nil {
			break // Nếu chưa tồn tại thì thoát khỏi vòng lặp
		}
	}

	// Tạo bản ghi mới trong DB
	url := models.URL{
		LongURL:  request.LongURL,
		ShortURL: short,
	}
	database.DB.Create(&url) // Lưu vào cơ sở dữ liệu

	// Trả về short URL cho client
	c.JSON(http.StatusOK, gin.H{"short_url": short})
}

// API GET /:code – Nhận short_url và chuyển hướng đến long_url
func RedirectURL(c *gin.Context) {
	code := c.Param("code") // dùng để lấy tham số từ URL
	var url models.URL      // biến kiểu models.URL để lưu kết quả truy vấn từ database.

	// Tìm bản ghi tương ứng với short code
	if err := database.DB.Where("short_url = ?", code).First(&url).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"}) // Nếu không thấy, trả về 404
		return
	}

	url.Clicks++           // Tăng số lần click lên 1
	database.DB.Save(&url) // Cập nhật lại DB

	c.Redirect(http.StatusMovedPermanently, url.LongURL) // Chuyển hướng 301 đến long URL
}

// API GET /stats/:code – Lấy thông tin thống kê lượt click của short_url
func GetStats(c *gin.Context) {
	code := c.Param("code") // Lấy short code từ URL path
	var url models.URL

	// Tìm bản ghi tương ứng trong DB
	if err := database.DB.Where("short_url = ?", code).First(&url).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"}) // Trả về lỗi nếu không tồn tại
		return
	}

	// Trả về thông tin thống kê
	c.JSON(http.StatusOK, gin.H{
		"long_url": url.LongURL,
		"clicks":   url.Clicks, // Số lần URL được click
	})
}
