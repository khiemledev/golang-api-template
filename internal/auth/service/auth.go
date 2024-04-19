package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"khiemle.dev/golang-api-template/internal/user/model"
	_userService "khiemle.dev/golang-api-template/internal/user/service"
)

type AuthService interface {
	LoginByUsernamePassword(ctx *gin.Context, username string, password string) (*model.User, error)
	RegisterUser(ctx *gin.Context, username string, email string, name string, password string, confirmPassword string) (*model.User, error)
}

type authService struct {
	db          *gorm.DB
	userService _userService.UserService
}

func NewAuthService(db *gorm.DB, userService _userService.UserService) AuthService {
	return &authService{
		db:          db,
		userService: userService,
	}
}

func (s *authService) LoginByUsernamePassword(ctx *gin.Context, username string, password string) (*model.User, error) {
	user, err := s.userService.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
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
