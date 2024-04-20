package token

import (
	"encoding/json"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/rs/zerolog/log"
	"khiemle.dev/golang-api-template/pkg/util"
)

type TokenMaker interface {
	GenerateToken(payload *TokenPayload, expDuration time.Duration) (string, error)
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

func (m *tokenMaker) GenerateToken(payload *TokenPayload, expDuration time.Duration) (string, error) {
	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(expDuration))

	// Set payload
	data, err := json.Marshal(payload)
	if err != nil {
		return "", nil
	}
	var json_data map[string]interface{}
	err = json.Unmarshal(data, &json_data)
	if err != nil {
		return "", err
	}

	for k, v := range json_data {
		err := token.Set(k, v)
		if err != nil {
			return "", err
		}
	}

	secretKeyHex := m.cfg.PasetoSecretHex
	key, _ := paseto.V4SymmetricKeyFromHex(secretKeyHex)

	signed := token.V4Encrypt(key, nil)
	return signed, nil
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
	err = json.Unmarshal(parsedToken.ClaimsJSON(), payload)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	return payload, nil
}
