package main

import (
	"github.com/VaheMuradyan/CodeSignal2/todoapp/middleware"
	"github.com/VaheMuradyan/CodeSignal2/todoapp/repositories/db"
	"github.com/VaheMuradyan/CodeSignal2/todoapp/router"
	"github.com/gin-gonic/gin"
)

func main() {
	database := db.ConnectDatabase()

	r := gin.Default()
	r.Use(middleware.RateLimiterMiddleware())

	router.RegisterRoutes(r, database)

	r.Run(":8080")
}
