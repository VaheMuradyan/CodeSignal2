package db

import (
	"time"

	"github.com/VaheMuradyan/CodeSignal2/todoapp/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	dsn := "root:java@tcp(127.0.0.1:3306)/todos?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("faild to connect database")
	}

	db.AutoMigrate(&models.Todo{}, &models.Library{})
	return db
}

func Reset(db *gorm.DB) {
	db.Exec("DROP TABLE IF EXISTS todos")
	db.AutoMigrate(&models.Todo{}) // This recreates the table

	// Seed with initial data
	db.Create(&[]models.Todo{
		{Title: "Learn Go", DueDate: time.Now().AddDate(0, 0, 5), Completed: false},
		{Title: "Build an API", DueDate: time.Now().AddDate(0, 0, 9), Completed: false},
		{Title: "Write Tests", DueDate: time.Now().AddDate(0, 0, 6), Completed: true},
	})
}
