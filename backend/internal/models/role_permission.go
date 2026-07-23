package models

import "time"

type RolePermission struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	RoleID       uint       `gorm:"not null;uniqueIndex:idx_role_permission" json:"role_id"`
	PermissionID uint       `gorm:"not null;uniqueIndex:idx_role_permission" json:"permission_id"`
	Role         Role       `json:"role"`
	Permission   Permission `json:"permission"`
	CreatedAt    time.Time  `json:"created_at"`
}
