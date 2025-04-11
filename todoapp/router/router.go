package router

import (
	"github.com/VaheMuradyan/CodeSignal2/todoapp/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/api/todos/due-soon", func(c *gin.Context) {
		controllers.GetTodosDueSoonHandler(c, db)
	})

	router.GET("/api/todos/overview", func(c *gin.Context) {
		controllers.GetOverviewHandler(c, db)
	})

	router.GET("/api/todos", func(c *gin.Context) {
		controllers.GetTodosHandler(c, db)
	})

	router.POST("/api/todos", func(c *gin.Context) {
		controllers.CreateTodoHandler(c, db)
	})

	router.DELETE("/api/reset", func(c *gin.Context) {
		controllers.ResetTodosHandler(c, db)
	})

	router.GET("/api/libraries", func(c *gin.Context) {
		controllers.GetFilteredLibrariesHandler(c, db)
	})
}
