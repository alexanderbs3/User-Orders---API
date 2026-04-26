package model

import "time"



type User struct {

	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Name      string    `gorm:"type:varchar(100);not null" json:"name"`
    Email     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`

	 Orders []Order `gorm:"foreignKey:UserID" json:"orders,omitempty"`
}

