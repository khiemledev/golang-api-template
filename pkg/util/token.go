package util

import (
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/rs/zerolog/log"
)

// GenerateToken generates a PASETO token with the given user ID and configuration.
//
// Parameters:
// - cfg: A pointer to the Config struct that contains the PASETO secret key.
// - userId: A string representing the user ID to be included in the token.
//
// Returns:
// - A string representing the signed PASETO token.
func GenerateToken(cfg *Config, userId string) string {
	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(24 * time.Hour))

	token.SetString("userId", userId)

	secretKeyHex := cfg.PasetoSecretHex
	key, _ := paseto.V4SymmetricKeyFromHex(secretKeyHex)

	signed := token.V4Encrypt(key, nil)
	return signed
}

// VerifyToken verifies a PASETO token with the given configuration and signed token.
//
// Parameters:
// - cfg: A pointer to the Config struct that contains the PASETO secret key.
// - signed: A string representing the signed PASETO token.
//
// Returns:
// - A boolean indicating whether the token is valid or not.
func VerifyToken(cfg *Config, signed string) bool {
	parser := paseto.NewParser()

	secretKeyHex := cfg.PasetoSecretHex
	key, _ := paseto.V4SymmetricKeyFromHex(secretKeyHex)

	parsedToken, err := parser.ParseV4Local(key, signed, nil)
	if err != nil {
		log.Error().Msg(err.Error())
		return false
	}

	// print value of payload in parsedToken
	userId, err := parsedToken.GetString("userId")
	if err != nil {
		log.Error().Msg("Error")
	}
	log.Info().Msgf("UserId: %s", userId)

	return true
}
