package service

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"khiemle.dev/golang-api-template/internal/user/model"
	_userService "khiemle.dev/golang-api-template/internal/user/service"
	"khiemle.dev/golang-api-template/pkg/util/token"
)

type AuthService interface {
	LoginByUsernamePassword(ctx *gin.Context, username string, password string) (*model.User, string, error)
	RegisterUser(ctx *gin.Context, username string, email string, name string, password string, confirmPassword string) (*model.User, error)
}

type authService struct {
	db          *gorm.DB
	userService _userService.UserService
	tokenMaker  token.TokenMaker
}

func NewAuthService(db *gorm.DB, userService _userService.UserService, tokenMaker token.TokenMaker) AuthService {
	return &authService{
		db:          db,
		userService: userService,
		tokenMaker:  tokenMaker,
	}
}

func (s *authService) LoginByUsernamePassword(ctx *gin.Context, username string, password string) (*model.User, string, error) {
	user, err := s.userService.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return nil, "", err
	}

	// Generate token
	payload := token.NewPayload(fmt.Sprint(user.ID))
	log.Info().Msg(payload.ID.String())
	token := s.tokenMaker.GenerateToken(payload)

	return user, token, nil
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
