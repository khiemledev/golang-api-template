package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"khiemle.dev/golang-api-template/internal/schemas"
	"khiemle.dev/golang-api-template/pkg/util/token"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

func AuthorizationMiddleware(tokenMaker token.TokenMaker) gin.HandlerFunc {
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

		ctx.Set(AuthorizationPayloadKey, *payload)
		ctx.Next()
	}
}
