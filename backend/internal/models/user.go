package models

import "time"

const (
	UserStatusEnabled  = "enabled"
	UserStatusDisabled = "disabled"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:64;uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	DisplayName  string    `gorm:"size:64" json:"display_name"`
	Status       string    `gorm:"size:20;not null;default:enabled" json:"status"`
	RoleID       uint      `gorm:"not null;index" json:"role_id"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
