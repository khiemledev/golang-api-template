package model

import "time"

type User struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Username       string    `json:"username" gorm:"size:255;uniqueIndex;not null"`
	Email          string    `json:"email" gorm:"size:255;uniqueIndex;not null"`
	Name           string    `json:"name" gorm:"size:255;not null"`
	HashedPassword string    `json:"hashed_password" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
