package routes

import (
	"github.com/gin-gonic/gin"
	"khiemle.dev/golang-api-template/internal/todo/handler"
)

func SetupTodoRouter(todoGroup *gin.RouterGroup, todoHandler handler.TodoHandler) {
	todoGroup.GET("/", todoHandler.ListTodoHandler)
	todoGroup.GET("/:id", todoHandler.GetByIdHandler)
	todoGroup.POST("/", todoHandler.CreateTodoHandler)
	todoGroup.PATCH("/:id", todoHandler.UpdateTodoHandler)
	todoGroup.DELETE("/:id", todoHandler.DeleteTodoHandler)
}
