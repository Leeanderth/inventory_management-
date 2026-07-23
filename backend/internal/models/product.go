package models

import "time"

type Product struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:128;not null;index" json:"name"`
	SKU       string    `gorm:"size:64;uniqueIndex;not null" json:"sku"`
	Category  string    `gorm:"size:64;index" json:"category"`
	Quantity  int       `gorm:"not null;default:0" json:"quantity"`
	Remark    string    `gorm:"type:text" json:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
