package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"khiemle.dev/golang-api-template/internal/auth/service"
	"khiemle.dev/golang-api-template/internal/schemas"
	_userService "khiemle.dev/golang-api-template/internal/user/service"
	"khiemle.dev/golang-api-template/pkg/util/token"
)

const (
	AuthorizationHeaderKey   = "authorization"
	AuthorizationHeaderToken = "authorization_token"
	AuthorizationTypeBearer  = "bearer"
	AuthorizationPayloadKey  = "authorization_payload"
	AuthorizationCurrentUser = "current_user"
)

func GetBearerTokenMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, schemas.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) != 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, schemas.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}

		authorizationHeaderType := strings.ToLower(fields[0])
		if authorizationHeaderType != AuthorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationHeaderType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, schemas.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}

		token := fields[1]
		ctx.Set(AuthorizationHeaderToken, token)
	}
}

func VerifyTokenMiddleware(tokenMaker token.TokenMaker, loginSessionService service.LoginSessionService, userService _userService.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.MustGet(AuthorizationHeaderToken).(string)

		// Verify token
		payload, err := tokenMaker.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, schemas.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}

		// Check login session
		loginSession, err := loginSessionService.FindByTokenID(ctx, payload.TokenID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, schemas.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: "login session not found",
				Data:    nil,
			})
			return
		}
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, schemas.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
		log.Info().Msgf("Login session found: %d", loginSession.ID)

		// Get current logged in user
		user, err := userService.GetUserById(ctx, loginSession.UserId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, schemas.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: "user not found",
				Data:    nil,
			})
			return
		}
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, schemas.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
		log.Info().Msgf("User found: %d", user.ID)

		ctx.Set(AuthorizationHeaderToken, token)
		ctx.Set(AuthorizationPayloadKey, *payload)
		ctx.Set(AuthorizationCurrentUser, user)
		ctx.Next()
	}
}
