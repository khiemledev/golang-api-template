package schemas

import "khiemle.dev/golang-api-template/internal/todo/model"

type CreateTodoRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type ListTodoResponse struct {
	Status  int          `json:"status" binding:"required"`
	Message string       `json:"message" binding:"required"`
	Todos   []model.Todo `json:"todos" binding:"required"`
}

type GetTodoByIdResponse struct {
	Status  int        `json:"status" binding:"required"`
	Message string     `json:"message" binding:"required"`
	Todo    model.Todo `json:"todo" binding:"required"`
}

type UpdateTodoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
