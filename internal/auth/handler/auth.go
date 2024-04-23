package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"khiemle.dev/golang-api-template/internal/auth/service"
	"khiemle.dev/golang-api-template/internal/constant"
	"khiemle.dev/golang-api-template/internal/schemas"
	"khiemle.dev/golang-api-template/internal/user/model"
	"khiemle.dev/golang-api-template/pkg/middleware"
	"khiemle.dev/golang-api-template/pkg/util"
	_token "khiemle.dev/golang-api-template/pkg/util/token"
)

type AuthHandler interface {
	LoginHandler(c *gin.Context)
	RegisterHandler(c *gin.Context)
	VerifyAccessToken(c *gin.Context)
	LogoutHandler(c *gin.Context)
	RefreshTokenHandler(c *gin.Context)
}

type authHandler struct {
	cfg                 *util.Config
	authService         service.AuthService
	loginSessionService service.LoginSessionService
}

func NewAuthHandler(cfg *util.Config, authService service.AuthService, loginSessionService service.LoginSessionService) AuthHandler {
	return &authHandler{
		cfg:                 cfg,
		authService:         authService,
		loginSessionService: loginSessionService,
	}
}

// LoginHandler godoc
//	@Summary		Login with username and password
//	@Description	Login with username and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		schemas.AuthLoginRequest	true	"Enter username and password"
//	@Success		200		{object}	schemas.AuthLoginResponse
//	@Failure		400		{object}	schemas.APIResponse
//	@Failure		401		{object}	schemas.APIResponse
//	@Failure		500		{object}	schemas.APIResponse
//	@Router			/auth/login [post]
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

	loginSession, err := h.loginSessionService.Create(
		c,
		data.Payload.TokenID,
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
		LoginSessionID:        loginSession.ID,
		AccessToken:           data.AccessToken,
		RefreshToken:          data.RefreshToken,
		AccessTokenExpiresIn:  data.AccessTokenExpiresIn,
		RefreshTokenExpiresIn: data.RefreshTokenExpiresIn,
	})
}

// RegisterHandler godoc
//	@Summary		Register new user
//	@Description	Register new user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		schemas.AuthRegisterRequest	true	"Enter user information"
//	@Success		200		{object}	schemas.AuthRegisterResponse
//	@Failure		400		{object}	schemas.APIResponse
//	@Failure		500		{object}	schemas.APIResponse
//	@Router			/auth/register [post]
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

// RegisterHandler godoc
//	@Summary		Verify access token
//	@Description	Verify access token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	schemas.AuthLoginUserResponse
//	@Failure		401	{object}	schemas.APIResponse
//	@Failure		500	{object}	schemas.APIResponse
//	@Router			/auth/verify_access_token [get]
//	@Security		BearerAuth
func (h *authHandler) VerifyAccessToken(c *gin.Context) {
	payload := c.MustGet(middleware.AuthorizationPayloadKey).(_token.TokenPayload)
	currentUser := c.MustGet(middleware.AuthorizationCurrentUser).(*model.User)

	c.JSON(http.StatusOK, schemas.APIResponse{
		Status:  http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data: gin.H{
			"payload": payload,
			"user": schemas.AuthLoginUserResponse{
				ID:       currentUser.ID,
				Name:     currentUser.Name,
				Username: currentUser.Username,
				Email:    currentUser.Email,
			},
		},
	})
}

// LogoutHandler godoc
//	@Summary		Logout
//	@Description	Logout
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	schemas.APIResponse
//	@Failure		401	{object}	schemas.APIResponse
//	@Failure		500	{object}	schemas.APIResponse
//	@Router			/auth/logout [get]
//	@Security		BearerAuth
func (h *authHandler) LogoutHandler(c *gin.Context) {
	payload := c.MustGet(middleware.AuthorizationPayloadKey).(_token.TokenPayload)

	// Delete login session
	err := h.loginSessionService.DeleteByTokenID(c, payload.TokenID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, schemas.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.SetCookie(constant.AccessToken, "", -1, "/", h.cfg.CookieDomain, false, true)
	c.SetCookie(constant.RefreshToken, "", -1, "/", h.cfg.CookieDomain, false, true)

	c.JSON(http.StatusOK, schemas.APIResponse{
		Status:  http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    nil,
	})
}

// LogoutHandler godoc
//	@Summary		Refresh token
//	@Description	Refresh token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	schemas.AuthRefreshResponse
//	@Failure		401	{object}	schemas.APIResponse
//	@Failure		500	{object}	schemas.APIResponse
//	@Router			/auth/refresh_token [get]
//	@Security		BearerAuth
func (h *authHandler) RefreshTokenHandler(ctx *gin.Context) {
	refreshToken := ctx.MustGet(middleware.AuthorizationHeaderToken).(string)

	user, data, err := h.authService.RefreshToken(ctx, refreshToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, schemas.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.SetCookie(constant.AccessToken, data.AccessToken, int(data.AccessTokenExpiresIn.Seconds()), "/", h.cfg.CookieDomain, false, true)

	ctx.JSON(http.StatusOK, schemas.AuthRefreshResponse{
		Status:  http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		User: schemas.AuthLoginUserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
		},
		LoginSessionID:       data.LoginSessionID,
		AccessToken:          data.AccessToken,
		AccessTokenExpiresIn: data.AccessTokenExpiresIn,
	})
}
