package model

import (
	"time"

	"github.com/google/uuid"
	"khiemle.dev/golang-api-template/internal/user/model"
)

type LoginSession struct {
	ID                    uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	TokenID               uuid.UUID  `json:"token_id" gorm:"not null;uniqueIndex"`
	UserId                uint       `json:"user_id" gorm:"not null"`
	User                  model.User `json:"user" gorm:"foreignKey:UserId;references:ID"`
	AccessToken           string     `json:"access_token" gorm:"not null"`
	RefreshToken          string     `json:"refresh_token" gorm:"not null"`
	UserAgent             string     `json:"user_agent" gorm:"not null"`
	ClientIP              string     `json:"client_ip" gorm:"not null"`
	AccessTokenExpiresIn  time.Time  `json:"access_token_expires_in" gorm:"not null"`
	RefreshTokenExpiresIn time.Time  `json:"refresh_token_expires_in" gorm:"not null"`
	LastUsedAt            *time.Time `json:"last_used_at"`
	CreatedAt             time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt             time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt             *time.Time `json:"deleted_at" gorm:"index"`
}
