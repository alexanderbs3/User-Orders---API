package model

import "time"


type OrderStatus string

const (
    StatusPending  OrderStatus = "pending"
    StatusPaid     OrderStatus = "paid"
    StatusCanceled OrderStatus = "canceled"
)



type Order struct {
    ID          uint        `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID      uint        `gorm:"not null;index" json:"user_id"`
    Description string      `gorm:"type:text;not null" json:"description"`
    Amount      float64     `gorm:"type:decimal(10,2);not null" json:"amount"`
    Status      OrderStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`

    // BelongsTo: cada pedido pertence a um usuário
    User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}