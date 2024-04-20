package token

import (
	"github.com/google/uuid"
)

type TokenPayload struct {
	ID     uuid.UUID `json:"id"`
	UserId string    `json:"user_id"`
}

func NewPayload(userId string) *TokenPayload {
	return &TokenPayload{
		ID:     uuid.New(),
		UserId: userId,
	}
}
