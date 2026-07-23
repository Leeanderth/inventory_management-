package models

import "time"

type StockMovement struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ProductID      uint      `gorm:"not null;index" json:"product_id"`
	Product        Product   `json:"product"`
	BeforeQuantity int       `gorm:"not null" json:"before_quantity"`
	AfterQuantity  int       `gorm:"not null" json:"after_quantity"`
	ChangeQuantity int       `gorm:"not null" json:"change_quantity"`
	OperatorID     uint      `gorm:"not null;index" json:"operator_id"`
	Operator       User      `gorm:"foreignKey:OperatorID" json:"operator"`
	Remark         string    `gorm:"type:text" json:"remark"`
	CreatedAt      time.Time `json:"created_at"`
}
