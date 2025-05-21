package router

import (
	"github.com/VaheMuradyan/CodeSignal2/todoapp/controllers"
	"github.com/VaheMuradyan/CodeSignal2/todoapp/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	router := r.Group("/api")
	{
		router.GET("/todos/due-soon", func(c *gin.Context) {
			controllers.GetTodosDueSoonHandler(c, db)
		})

		router.GET("/todos/overview", func(c *gin.Context) {
			controllers.GetOverviewHandler(c, db)
		})

		router.DELETE("/reset", func(c *gin.Context) {
			controllers.ResetTodosHandler(c, db)
		})

		router.GET("/libraries", func(c *gin.Context) {
			controllers.GetFilteredLibrariesHandler(c, db)
		})

		router.GET("/ws", controllers.WebSocketHandler)

		router.GET("/todos/:id", controllers.GetTodoHandler(db))

		router.POST("/register", func(c *gin.Context) { controllers.Register(c, db) })
		router.POST("/login", func(c *gin.Context) { controllers.LoginWithSession(c) })

		protected := r.Group("/api", middleware.JWTAuthMiddleware())
		{
			protected.GET("/todos", func(c *gin.Context) {
				controllers.GetTodosHandler(c, db)
			})

			protected.POST("/todos", func(c *gin.Context) {
				controllers.CreateTodoHandler(c, db)
			})

		}
	}
}
