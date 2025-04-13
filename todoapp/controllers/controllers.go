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

func WebSocketHandler(c *gin.Context) {
	err := services.CheckConnectionCount()
	if err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		return
	}

	conn, err := services.UpgradeConnection(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to establish WebSocket connection"})
		return
	}
	defer services.CloseConnection(conn)

	services.WebSocketLoop(conn)
}

func GetTodoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Handler logic here
		id := c.Param("id")
		todo, err := services.GetTodoService(db, id)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "faild to get data"})
			return
		}
		c.JSON(http.StatusOK, todo)
	}
}

// =================== User ======================

func Register(c *gin.Context, db *gorm.DB) {
	var temp struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := services.RegisterUser(db, temp.Username, temp.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

func Login(c *gin.Context, db *gorm.DB) {
	var temp struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := services.ValidateUserCredentials(db, temp.Username, temp.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": "Login successful"})
}
