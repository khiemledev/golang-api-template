package service

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"khiemle.dev/golang-api-template/internal/user/model"
	_userService "khiemle.dev/golang-api-template/internal/user/service"
	"khiemle.dev/golang-api-template/pkg/util"
	"khiemle.dev/golang-api-template/pkg/util/token"
)

type LoginByUserNamePasswordData struct {
	Payload               *token.TokenPayload
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresIn  time.Duration
	RefreshTokenExpiresIn time.Duration
}

type RefreshTokenData struct {
	LoginSessionID       uint
	Payload              *token.TokenPayload
	AccessToken          string
	AccessTokenExpiresIn time.Duration
}

type AuthService interface {
	LoginByUsernamePassword(ctx *gin.Context, username string, password string) (*model.User, *LoginByUserNamePasswordData, error)
	RegisterUser(ctx *gin.Context, username string, email string, name string, password string, confirmPassword string) (*model.User, error)
	RefreshToken(ctx *gin.Context, refreshToken string) (*model.User, *RefreshTokenData, error)
}

type authService struct {
	db                  *gorm.DB
	cfg                 *util.Config
	userService         _userService.UserService
	loginSessionService LoginSessionService
	tokenMaker          token.TokenMaker
}

func NewAuthService(
	db *gorm.DB,
	cfg *util.Config,
	userService _userService.UserService,
	loginSessionService LoginSessionService,
	tokenMaker token.TokenMaker) AuthService {
	return &authService{
		db:                  db,
		cfg:                 cfg,
		userService:         userService,
		loginSessionService: loginSessionService,
		tokenMaker:          tokenMaker,
	}
}

func (s *authService) LoginByUsernamePassword(ctx *gin.Context, username string, password string) (*model.User, *LoginByUserNamePasswordData, error) {
	user, err := s.userService.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return nil, nil, err
	}

	// Generate token
	payload := token.NewPayload(user.ID)
	accessToken, err := s.tokenMaker.GenerateToken(payload, time.Duration(s.cfg.AccessTokenExpiryInHours)*time.Hour)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, err := s.tokenMaker.GenerateToken(payload, time.Duration(s.cfg.RefreshTokenExpiryInHours)*time.Hour)
	if err != nil {
		return nil, nil, err
	}

	data := &LoginByUserNamePasswordData{
		Payload:               payload,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresIn:  time.Duration(s.cfg.AccessTokenExpiryInHours) * time.Hour,
		RefreshTokenExpiresIn: time.Duration(s.cfg.RefreshTokenExpiryInHours) * time.Hour,
	}

	return user, data, nil
}

func (s *authService) RegisterUser(ctx *gin.Context, username string, email string, name string, password string, confirmPassword string) (*model.User, error) {
	if password != confirmPassword {
		return nil, errors.New("password and confirm password does not match")
	}

	user, err := s.userService.CreateUser(ctx, username, email, name, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) RefreshToken(ctx *gin.Context, refreshToken string) (*model.User, *RefreshTokenData, error) {
	payload, err := s.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		return nil, nil, err
	}

	loginSession, err := s.loginSessionService.FindByTokenID(ctx, payload.TokenID)
	if err != nil {
		return nil, nil, err
	}
	if refreshToken != loginSession.RefreshToken {
		return nil, nil, errors.New("refreshToken does not match")
	}

	user, err := s.userService.GetUserById(ctx, loginSession.UserId)
	if err != nil {
		return nil, nil, err
	}

	exp := time.Duration(s.cfg.AccessTokenExpiryInHours) * time.Hour
	accessToken, err := s.tokenMaker.GenerateToken(payload, exp)
	if err != nil {
		return nil, nil, err
	}

	loginSession, err = s.loginSessionService.UpdateAccessToken(ctx, loginSession.ID, accessToken, time.Now().Add(exp))
	if err != nil {
		return nil, nil, err
	}

	return user, &RefreshTokenData{
		LoginSessionID:       loginSession.ID,
		Payload:              payload,
		AccessToken:          accessToken,
		AccessTokenExpiresIn: exp,
	}, nil
}
