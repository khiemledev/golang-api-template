package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"khiemle.dev/golang-api-template/internal/schemas"
	"khiemle.dev/golang-api-template/internal/todo/service"
)

type TodoHandler interface {
	CreateTodoHandler(c *gin.Context)
	GetByIdHandler(c *gin.Context)
	UpdateTodoHandler(c *gin.Context)
	ListTodoHandler(c *gin.Context)
	DeleteTodoHandler(c *gin.Context)
}

type todoHandler struct {
	todoService service.TodoService
}

func NewTodoHandler(todoService service.TodoService) TodoHandler {
	return &todoHandler{
		todoService: todoService,
	}
}

// Handler for create todo
func (h *todoHandler) CreateTodoHandler(c *gin.Context) {
	log.Info().Msg("I'm in")
	req := schemas.CreateTodoRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, schemas.APIResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	todo, err := h.todoService.CreateTodo(c, req.Name, req.Description)
	if err != nil {
		log.Fatal().Msg("Error creating todo")
		log.Error().Err(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, schemas.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, schemas.APIResponse{
		Status:  http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    todo,
	})
}

// Handler for get todo
func (h *todoHandler) GetByIdHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, schemas.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid id",
			Data:    nil,
		})
		return
	}

	todo, err := h.todoService.GetById(c, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, schemas.APIResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, schemas.GetTodoByIdResponse{
		Status:  http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Todo:    *todo,
	})
}

// Handler for update todo
func (h *todoHandler) UpdateTodoHandler(c *gin.Context) {
	req := schemas.UpdateTodoRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, schemas.APIResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, schemas.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid id",
			Data:    nil,
		})
		return
	}

	todo, err := h.todoService.UpdateTodo(c, id, req.Name, req.Description)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, schemas.APIResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
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

	c.JSON(http.StatusOK, schemas.APIResponse{
		Status:  http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    todo,
	})
}

// Handler for list todo
func (h *todoHandler) ListTodoHandler(c *gin.Context) {
	todos := h.todoService.ListTodo(c)

	c.JSON(http.StatusOK, schemas.ListTodoResponse{
		Status:  http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Todos:   todos,
	})
}

// Handler for delete todo
func (h *todoHandler) DeleteTodoHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, schemas.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid id",
			Data:    nil,
		})
		return
	}

	err = h.todoService.DeleteTodo(c, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, schemas.APIResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
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

	c.JSON(http.StatusOK, schemas.APIResponse{
		Status:  http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    nil,
	})
}
