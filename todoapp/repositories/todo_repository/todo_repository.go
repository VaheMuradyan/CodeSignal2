package todo_repository

import (
	"time"

	"github.com/VaheMuradyan/CodeSignal2/todoapp/models"
	"gorm.io/gorm"
)

func FindTodosDueSoon(db *gorm.DB, nextWeek time.Time) []models.Todo {
	var todos []models.Todo
	db.Where("due_date <= ? AND completed = ? AND priority = ?", nextWeek, false, "high").Find(&todos)
	return todos
}

func GetTodoOverview(db *gorm.DB) (int64, int64) {
	var total, completed int64
	db.Model(&models.Todo{}).Count(&total)
	db.Model(&models.Todo{}).Where("completed = ?", true).Count(&completed)
	return total, completed
}

func FindAllTodos(db *gorm.DB) []models.Todo {
	var todo []models.Todo
	db.Find(&todo)
	return todo
}

func CreateTodo(db *gorm.DB, todo models.Todo) models.Todo {
	db.Create(&todo)
	return todo
}

func ResetTodos(db *gorm.DB) {
	db.Exec("DELETE FROM todos")
	db.Exec("ALTER TABLE todos AUTO_INCREMENT = 1")
}

func GetFilteredLibraries(db *gorm.DB, location string, locationExists bool, libraryType string, typeExists bool, openStatus string, openStatusExists bool) []models.Library {
	var libraries []models.Library

	// TODO: Initialize the query
	query := db.Model(&models.Library{})
	// TODO: add conditions dynamically based on the presence of query parameters
	if locationExists {
		query = query.Where("location = ?", location)
	}
	if typeExists {
		query = query.Where("type = ?", libraryType)
	}
	if openStatusExists {
		query = query.Where("is_open = ?", openStatus == "true")
	}
	// TODO: execute query
	query.Find(&libraries)

	return libraries
}
