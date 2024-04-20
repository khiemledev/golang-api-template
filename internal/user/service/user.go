package service

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"khiemle.dev/golang-api-template/internal/user/model"
)

type UserService interface {
	CreateUser(ctx *gin.Context, username string, email string, name string, password string) (*model.User, error)
	GetUserById(ctx *gin.Context, id uint) (*model.User, error)
	GetUserByUsername(ctx *gin.Context, username string) (*model.User, error)
	GetUserByEmail(ctx *gin.Context, email string) (*model.User, error)
	ListUser(ctx *gin.Context) ([]model.User, error)
	UpdateUser(ctx *gin.Context, id uint, name string, email string) (*model.User, error)
	UpdatePassword(ctx *gin.Context, id uint, password string, newPassword string, confirmPassword string) error
	DeleteUser(ctx *gin.Context, id uint) error
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{
		db: db,
	}
}

func (s *userService) CreateUser(ctx *gin.Context, username string, email string, name string, password string) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:       username,
		Email:          email,
		Name:           name,
		HashedPassword: string(hashedPassword),
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&user).Error
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserById(ctx *gin.Context, id uint) (*model.User, error) {
	user := &model.User{}
	tx := s.db.First(&user, "id = ?", id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}

func (s *userService) GetUserByUsername(ctx *gin.Context, username string) (*model.User, error) {
	user := &model.User{}
	tx := s.db.First(&user, "username = ?", username)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}

func (s *userService) GetUserByEmail(ctx *gin.Context, email string) (*model.User, error) {
	user := &model.User{}
	tx := s.db.First(&user, "email = ?", email)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}

func (s *userService) ListUser(ctx *gin.Context) ([]model.User, error) {
	users := []model.User{}
	tx := s.db.Find(&users)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return users, nil
}

func (s *userService) UpdateUser(ctx *gin.Context, id uint, name string, email string) (*model.User, error) {
	user := &model.User{}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Find(user, "id = ?", id).Error
		if err != nil {
			return err
		}

		if name != "" {
			user.Name = name
		}
		if email != "" {
			user.Email = email
		}

		return tx.Save(user).Error
	})
	return user, err
}

func (s *userService) UpdatePassword(ctx *gin.Context, id uint, currentPassword string, newPassword string, confirmPassword string) error {
	if newPassword != confirmPassword {
		return nil
	}

	user := &model.User{}
	err := s.db.First(user, "id = ?", id).Error
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(currentPassword)); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		user.HashedPassword = string(hashedPassword)
		return tx.Save(user).Error
	})
	return err
}

func (s *userService) DeleteUser(ctx *gin.Context, id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.User{}, "id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}
