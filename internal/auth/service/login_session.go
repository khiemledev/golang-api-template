package service

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"khiemle.dev/golang-api-template/internal/auth/model"
)

type LoginSessionService interface {
	Create(c *gin.Context, tokenID uuid.UUID, userId uint, userAgent string, clientIP string, accessToken string, refreshToken string, accessTokenExpiresIn time.Time, refreshTokenExpiresIn time.Time) (*model.LoginSession, error)
	FindById(c *gin.Context, id uint) (*model.LoginSession, error)
	FindByTokenID(c *gin.Context, tokenID uuid.UUID) (*model.LoginSession, error)
	UpdateAccessToken(c *gin.Context, id uint, accessToken string, accessTokenExpiresIn time.Time) (*model.LoginSession, error)
	DeleteByTokenID(c *gin.Context, tokenID uuid.UUID) error
}

type loginSessionService struct {
	db *gorm.DB
}

func NewLoginSessionService(db *gorm.DB) LoginSessionService {
	return &loginSessionService{
		db: db,
	}
}

func (s *loginSessionService) FindById(c *gin.Context, id uint) (*model.LoginSession, error) {
	loginSession := model.LoginSession{}
	tx := s.db.First(&loginSession, "id = ?", id)
	return &loginSession, tx.Error
}

func (s *loginSessionService) FindByTokenID(c *gin.Context, tokenID uuid.UUID) (*model.LoginSession, error) {
	loginSession := model.LoginSession{}
	tx := s.db.First(&loginSession, "token_id = ?", tokenID)
	return &loginSession, tx.Error
}

func (s *loginSessionService) Create(c *gin.Context, tokenID uuid.UUID, userId uint, userAgent string, clientIP string, accessToken string, refreshToken string, accessTokenExpiresIn time.Time, refreshTokenExpiresIn time.Time) (*model.LoginSession, error) {
	loginSession := model.LoginSession{
		TokenID:               tokenID,
		UserId:                userId,
		UserAgent:             userAgent,
		ClientIP:              clientIP,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresIn:  accessTokenExpiresIn,
		RefreshTokenExpiresIn: refreshTokenExpiresIn,
	}
	tx := s.db.Create(&loginSession)
	return &loginSession, tx.Error
}

func (s *loginSessionService) UpdateAccessToken(c *gin.Context, id uint, accessToken string, accessTokenExpiresIn time.Time) (*model.LoginSession, error) {
	loginSession := model.LoginSession{}
	tx := s.db.Model(&loginSession).
		Where("id = ?", id).
		Update("access_token", accessToken).
		Update("access_token_expires_in", accessTokenExpiresIn)
	return &loginSession, tx.Error
}

func (s *loginSessionService) DeleteByTokenID(c *gin.Context, tokenID uuid.UUID) error {
	loginSession := model.LoginSession{}
	tx := s.db.Delete(&loginSession, "token_id = ?", tokenID)
	return tx.Error
}
