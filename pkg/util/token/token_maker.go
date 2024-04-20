package token

import (
	"encoding/json"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/rs/zerolog/log"
	"khiemle.dev/golang-api-template/pkg/util"
)

type TokenMaker interface {
	GenerateToken(payload *TokenPayload) string
	VerifyToken(signed string) (*TokenPayload, error)
}

type tokenMaker struct {
	cfg *util.Config
}

func NewTokenMaker(cfg *util.Config) TokenMaker {
	return &tokenMaker{
		cfg: cfg,
	}
}

func (m *tokenMaker) GenerateToken(payload *TokenPayload) string {
	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(30 * time.Second))

	// Set payload
	data, _ := json.Marshal(payload)
	var json_data map[string]interface{}
	json.Unmarshal(data, &json_data)
	for k, v := range json_data {
		token.Set(k, v)
		log.Info().Msg(k + " : " + v.(string))
	}

	secretKeyHex := m.cfg.PasetoSecretHex
	key, _ := paseto.V4SymmetricKeyFromHex(secretKeyHex)

	signed := token.V4Encrypt(key, nil)
	return signed
}

func (m *tokenMaker) VerifyToken(signed string) (*TokenPayload, error) {
	parser := paseto.NewParser()

	secretKeyHex := m.cfg.PasetoSecretHex
	key, _ := paseto.V4SymmetricKeyFromHex(secretKeyHex)

	parsedToken, err := parser.ParseV4Local(key, signed, nil)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	payload := &TokenPayload{}
	json.Unmarshal(parsedToken.ClaimsJSON(), payload)

	return payload, nil
}
