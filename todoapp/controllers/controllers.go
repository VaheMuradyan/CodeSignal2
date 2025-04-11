package controllers

import (
	"net/http"

	"github.com/VaheMuradyan/CodeSignal2/todoapp/models"
	"github.com/VaheMuradyan/CodeSignal2/todoapp/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTodosDueSoonHandler(c *gin.Context, db *gorm.DB) {
	todos := services.GetTodosDueSoon(db)
	c.JSON(http.StatusOK, todos)
}

func GetOverviewHandler(c *gin.Context, db *gorm.DB) {
	overview := services.GetOverview(db)
	c.JSON(http.StatusOK, overview)
}

func GetTodosHandler(c *gin.Context, db *gorm.DB) {
	todos := services.GetTodos(db)
	c.JSON(http.StatusOK, todos)
}

func CreateTodoHandler(c *gin.Context, db *gorm.DB) {
	var newTodo models.Todo
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	createdTodo := services.AddTodo(db, newTodo)
	c.JSON(http.StatusCreated, createdTodo)
}

func ResetTodosHandler(c *gin.Context, db *gorm.DB) {
	services.ResetAllTodos(db)
	c.Status(http.StatusOK)
}

func GetFilteredLibrariesHandler(c *gin.Context, db *gorm.DB) {
	libraries := services.GetFilteredLibrariesService(c, db)
	c.JSON(http.StatusOK, libraries)
}
