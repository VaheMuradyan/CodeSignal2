package services

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/VaheMuradyan/CodeSignal2/todoapp/models"
	"github.com/VaheMuradyan/CodeSignal2/todoapp/repositories/todo_repository"
	"github.com/VaheMuradyan/CodeSignal2/todoapp/repositories/user_repository"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
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

const maxConnections = 3

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return r.Host == "localhost:3000" || r.Host == "127.0.0.1:3000"
		},
	}

	connections   = make(map[*websocket.Conn]bool)
	connectionsMu sync.Mutex
)

func CheckConnectionCount() error {
	connectionsMu.Lock()
	defer connectionsMu.Unlock()

	if len(connections) >= maxConnections {
		return errors.New("to many connections")
	}
	return nil
}

func UpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {

	if err := CheckConnectionCount(); err != nil {
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return nil, err
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	connectionsMu.Lock()
	connections[conn] = true
	connectionsMu.Unlock()
	return conn, nil
}

func WebSocketLoop(conn *websocket.Conn) {
	defer CloseConnection(conn)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("[SERVER LOG] Unexpected WebSocket error:", err)
			}
			break
		}
		fmt.Println("[SERVER LOG] Message received: ", string(msg))
		for conn := range connections {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				fmt.Println("[SERVER LOG] Error writing message:", err)
				break
			}
		}
	}
}

func CloseConnection(conn *websocket.Conn) {
	connectionsMu.Lock()
	if _, ok := connections[conn]; ok {
		conn.Close()
		delete(connections, conn)
	}
	connectionsMu.Unlock()
}

func Reset() {
	connectionsMu.Lock()
	for conn := range connections {
		conn.Close()
		delete(connections, conn)
	}
	connectionsMu.Unlock()
}

func GetTodoService(db *gorm.DB, id string) (models.Todo, error) {
	return todo_repository.GetTodoByID(db, id)
}

// =================================== User Parts ==================================
func RegisterUser(db *gorm.DB, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Username: username,
		Password: string(hashedPassword),
	}

	return user_repository.CreateUser(db, &user)
}

func ValidateUserCredentials(db *gorm.DB, username, password string) error {
	user, err := user_repository.GetUserByUsername(db, username)
	if err != nil {
		return err
	}

	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
