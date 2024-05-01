package service

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"khiemle.dev/golang-api-template/internal/auth/model"
)

type CreateLoginSessionArgs struct {
	TokenID               uuid.UUID
	UserId                uint
	UserAgent             string
	ClientIP              string
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresIn  time.Time
	RefreshTokenExpiresIn time.Time
}

type UpdateLoginSessionArgs struct {
	AccessToken           *string
	RefreshToken          *string
	AccessTokenExpiresIn  *time.Time
	RefreshTokenExpiresIn *time.Time
	LastUsedAt            *time.Time
}

type LoginSessionService interface {
	Create(c *gin.Context, args CreateLoginSessionArgs) (*model.LoginSession, error)
	FindById(c *gin.Context, id uint) (*model.LoginSession, error)
	FindByTokenID(c *gin.Context, tokenID uuid.UUID) (*model.LoginSession, error)
	UpdateSession(c *gin.Context, id uint, args UpdateLoginSessionArgs) (*model.LoginSession, error)
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

func (s *loginSessionService) Create(c *gin.Context, args CreateLoginSessionArgs) (*model.LoginSession, error) {
	loginSession := model.LoginSession{
		TokenID:               args.TokenID,
		UserId:                args.UserId,
		UserAgent:             args.UserAgent,
		ClientIP:              args.ClientIP,
		AccessToken:           args.AccessToken,
		RefreshToken:          args.RefreshToken,
		AccessTokenExpiresIn:  args.AccessTokenExpiresIn,
		RefreshTokenExpiresIn: args.RefreshTokenExpiresIn,
	}
	tx := s.db.Create(&loginSession)
	return &loginSession, tx.Error
}

func (s *loginSessionService) UpdateSession(c *gin.Context, id uint, args UpdateLoginSessionArgs) (*model.LoginSession, error) {
	loginSession := model.LoginSession{}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		err := tx.First(&loginSession, "id = ?", id).Error
		if err != nil {
			return err
		}

		if args.AccessToken != nil {
			loginSession.AccessToken = *args.AccessToken
		}
		if args.RefreshToken != nil {
			loginSession.RefreshToken = *args.RefreshToken
		}
		if args.AccessTokenExpiresIn != nil {
			loginSession.AccessTokenExpiresIn = *args.AccessTokenExpiresIn
		}
		if args.RefreshTokenExpiresIn != nil {
			loginSession.RefreshTokenExpiresIn = *args.RefreshTokenExpiresIn
		}
		if args.LastUsedAt != nil {
			loginSession.LastUsedAt = args.LastUsedAt
		}
		return tx.Save(&loginSession).Error
	})
	return &loginSession, err
}

func (s *loginSessionService) DeleteByTokenID(c *gin.Context, tokenID uuid.UUID) error {
	loginSession := model.LoginSession{}
	tx := s.db.Delete(&loginSession, "token_id = ?", tokenID)
	return tx.Error
}
