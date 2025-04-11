package services

import (
	"time"

	"github.com/VaheMuradyan/CodeSignal2/todoapp/models"
	"github.com/VaheMuradyan/CodeSignal2/todoapp/repositories/todo_repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTodosDueSoon(db *gorm.DB) []models.Todo {
	nextWeek := time.Now().AddDate(0, 0, 7)
	return todo_repository.FindTodosDueSoon(db, nextWeek)
}

func GetOverview(db *gorm.DB) map[string]int64 {
	total, completed := todo_repository.GetTodoOverview(db)
	return map[string]int64{
		"totalTodos": total,
		"completed":  completed,
		"incomplete": total - completed,
	}
}

func GetTodos(db *gorm.DB) []models.Todo {
	return todo_repository.FindAllTodos(db)
}

func AddTodo(db *gorm.DB, newTodo models.Todo) models.Todo {
	newTodo.Completed = false
	return todo_repository.CreateTodo(db, newTodo)
}

func ResetAllTodos(db *gorm.DB) {
	todo_repository.ResetTodos(db)
}

func GetFilteredLibrariesService(c *gin.Context, db *gorm.DB) []models.Library {
	location, locationExists := c.GetQuery("location")
	libraryType, typeExists := c.GetQuery("type")
	openStatus, openStatusExists := c.GetQuery("isOpen")

	return todo_repository.GetFilteredLibraries(db, location, locationExists, libraryType, typeExists, openStatus, openStatusExists)
}
