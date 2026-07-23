package models

import "time"

type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"size:64;not null" json:"name"`
	Description string    `gorm:"size:255" json:"description"`
	IsSystem    bool      `gorm:"not null;default:false" json:"is_system"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
