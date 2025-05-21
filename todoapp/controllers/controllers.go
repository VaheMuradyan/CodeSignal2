package controllers

import (
	"net/http"
	"time"

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
	var creds models.Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := services.RegisterUser(db, creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

func LoginWithSession(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if !services.Authenticate(credentials.Username, credentials.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	tokenString, err := services.CreateToken(credentials.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.SetCookie("session_token", tokenString, int(4*time.Hour), "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful, session created"})
}
func CheckSession(c *gin.Context) {
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session not found"})
		return
	}

	claims, err := services.ValidateToken(sessionToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session is valid", "username": claims.Username})
}

func GenerateToken(c *gin.Context) {
	token, err := services.CreateToken("testuser")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": token})
}
