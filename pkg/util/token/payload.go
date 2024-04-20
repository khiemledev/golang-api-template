package token

import (
	"github.com/google/uuid"
)

type TokenPayload struct {
	TokenID uuid.UUID `json:"id"`
	UserId  uint      `json:"user_id"`
}

func NewPayload(userId uint) *TokenPayload {
	return &TokenPayload{
		TokenID: uuid.New(),
		UserId:  userId,
	}
}
