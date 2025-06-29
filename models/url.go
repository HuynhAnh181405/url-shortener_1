package models

import "gorm.io/gorm"

// URL là struct đại diện cho bảng lưu trữ thông tin URL gốc và URL rút gọn
type URL struct {
	gorm.Model
	LongURL  string `gorm:"not null"`
	ShortURL string `gorm:"uniqueIndex"`
	Clicks   uint
}
