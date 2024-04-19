package service

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"khiemle.dev/golang-api-template/internal/todo/model"
)

type TodoService interface {
	CreateTodo(c *gin.Context, name string, description string) (*model.Todo, error)
	GetById(c *gin.Context, id int) (*model.Todo, error)
	UpdateTodo(c *gin.Context, id int, name string, description string) (*model.Todo, error)
	ListTodo(c *gin.Context) []model.Todo
	DeleteTodo(c *gin.Context, id int) error
}

type todoService struct {
	db *gorm.DB
}

func NewTodoService(db *gorm.DB) TodoService {
	return &todoService{
		db: db,
	}
}

// CreateTodo
func (s *todoService) CreateTodo(c *gin.Context, name string, description string) (*model.Todo, error) {
	todo := model.Todo{
		Name:        name,
		Description: description,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&todo).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

// GetById
func (s *todoService) GetById(c *gin.Context, id int) (*model.Todo, error) {
	todo := model.Todo{}
	tx := s.db.First(&todo, "id = ?", id)
	return &todo, tx.Error
}

// UpdateTodo
func (s *todoService) UpdateTodo(c *gin.Context, id int, name string, description string) (*model.Todo, error) {
	todo := model.Todo{}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		err := s.db.First(&todo, "id = ?", id).Error
		if err != nil {
			return err
		}

		if name != "" {
			todo.Name = name
		}
		if description != "" {
			todo.Description = description
		}

		if err := tx.Save(&todo).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

// ListTodo
func (s *todoService) ListTodo(c *gin.Context) []model.Todo {
	log.Info().Msg("I'm ok")
	todos := []model.Todo{}
	s.db.Find(&todos)

	return todos
}

// DeleteTodo
func (s *todoService) DeleteTodo(c *gin.Context, id int) error {
	todo := model.Todo{}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		err := s.db.First(&todo, "id = ?", id).Error
		if err != nil {
			return err
		}

		if err := tx.Delete(&todo).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
