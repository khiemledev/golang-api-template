package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"khiemle.dev/golang-api-template/internal/auth/service"
	"khiemle.dev/golang-api-template/internal/schemas"
	"khiemle.dev/golang-api-template/pkg/middleware"
	"khiemle.dev/golang-api-template/pkg/util"
	"khiemle.dev/golang-api-template/pkg/util/token"
)

type AuthHandler interface {
	LoginHandler(c *gin.Context)
	RegisterHandler(c *gin.Context)
	VerifyAccessToken(c *gin.Context)
}

type authHandler struct {
	cfg          *util.Config
	authService  service.AuthService
	loginService service.LoginSessionService
}

func NewAuthHandler(cfg *util.Config, authService service.AuthService, loginService service.LoginSessionService) AuthHandler {
	return &authHandler{
		cfg:          cfg,
		authService:  authService,
		loginService: loginService,
	}
}

func (h *authHandler) LoginHandler(c *gin.Context) {
	req := schemas.AuthLoginRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, schemas.APIResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	user, data, err := h.authService.LoginByUsernamePassword(c, req.Username, req.Password)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, schemas.APIResponse{
			Status:  http.StatusNotFound,
			Message: "username not found",
			Data:    nil,
		})
		return
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, schemas.APIResponse{
			Status:  http.StatusUnauthorized,
			Message: "Wrong password",
			Data:    nil,
		})
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, schemas.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	loginSession, err := h.loginService.Create(
		c,
		data.Payload.ID,
		user.ID,
		c.Request.UserAgent(),
		c.ClientIP(),
		data.AccessToken,
		data.RefreshToken,
		time.Now().Add(data.AccessTokenExpiresIn),
		time.Now().Add(data.RefreshTokenExpiresIn),
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, schemas.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.SetCookie("access_token", data.AccessToken, int(data.AccessTokenExpiresIn.Seconds()), "/", h.cfg.CookieDomain, false, true)
	c.SetCookie("refresh_token", data.RefreshToken, int(data.RefreshTokenExpiresIn.Seconds()), "/", h.cfg.CookieDomain, false, true)

	c.JSON(http.StatusOK, schemas.AuthLoginResponse{
		Status:  http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		User: schemas.AuthLoginUserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
		},
		LoginSessionID: loginSession.ID,
		AccessToken:    data.AccessToken,
		RefreshToken:   data.RefreshToken,
	})
}

func (h *authHandler) RegisterHandler(c *gin.Context) {
	req := schemas.AuthRegisterRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, schemas.APIResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// Validate email
	if !util.IsValidEmail(req.Email) {
		c.AbortWithStatusJSON(http.StatusBadRequest, schemas.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid email",
			Data:    nil,
		})
		return
	}

	// Validate password
	if !util.IsValidPassword(req.Password) {
		c.AbortWithStatusJSON(http.StatusBadRequest, schemas.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "Password must be at least 8 characters long and contain at least 1 uppercase letter, 1 lowercase letter, and 1 number",
			Data:    nil,
		})
		return
	}

	user, err := h.authService.RegisterUser(c, req.Username, req.Email, req.Name, req.Password, req.ConfirmPassword)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, schemas.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, schemas.AuthRegisterResponse{
		Status:        http.StatusOK,
		Message:       http.StatusText(http.StatusOK),
		CreatedUserId: user.ID,
	})
}

func (h *authHandler) VerifyAccessToken(c *gin.Context) {
	payload := c.MustGet(middleware.AuthorizationPayloadKey).(token.TokenPayload)

	c.JSON(http.StatusOK, schemas.APIResponse{
		Status:  http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    payload,
	})
}
