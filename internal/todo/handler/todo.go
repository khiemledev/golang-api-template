package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListAllTodos(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "List all todos",
	})
}
