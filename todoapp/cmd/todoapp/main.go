package main

import (
	"github.com/VaheMuradyan/CodeSignal2/todoapp/repositories/db"
	"github.com/VaheMuradyan/CodeSignal2/todoapp/router"
	"github.com/gin-gonic/gin"
)

func main() {
	database := db.ConnectDatabase()

	r := gin.Default()
	router.RegisterRoutes(r, database)

	// Start the server on port 8080
	r.Run(":8080")
}
