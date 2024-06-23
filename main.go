package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// todo represents data about a task.
type todo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// todos slice to seed record todo data.
var todos = []todo{
	{ID: "1", Title: "Learn Go", Completed: false},
	{ID: "2", Title: "Read Gin Documentation", Completed: false},
	{ID: "3", Title: "Build a REST API", Completed: false},
}
func getTodos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, todos)
}

func serveIndex(c *gin.Context) {
	c.File("./static/index.html")
}

var idCounter = 3
// var idMutex sync.Mutex

func postTodo(c *gin.Context) {
	var newTodo todo

	// Call BindJSON to bind the received JSON to newTodo.
	if err := c.BindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate input
	if newTodo.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	// Generate a new ID
	// idMutex.Lock()
	idCounter++
	newTodo.ID = strconv.Itoa(idCounter)
	// idMutex.Unlock()

	// Add the new todo to the slice.
	todos = append(todos, newTodo)
	c.IndentedJSON(http.StatusCreated, newTodo)
}

func getTodoByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of todos, looking for a todo whose ID value matches the parameter.
	for _, t := range todos {
		if t.ID == id {
			c.IndentedJSON(http.StatusOK, t)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
}

func updateTodo(c *gin.Context) {
	id := c.Param("id")
	var updatedTodo todo

	if err := c.BindJSON(&updatedTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate input
	if updatedTodo.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	for i, t := range todos {
		if t.ID == id {
			updatedTodo.ID = id // Ensure the ID remains the same
			todos[i] = updatedTodo
			c.IndentedJSON(http.StatusOK, updatedTodo)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
}

func deleteTodo(c *gin.Context) {
	id := c.Param("id")

	for i, t := range todos {
		if t.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "todo deleted"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
}

func main() {
	router := gin.New()
	gin.SetMode(gin.ReleaseMode)
	router.GET("/", serveIndex)
	router.GET("/todos", getTodos)
	router.POST("/todos", postTodo)
	router.GET("/todos/:id", getTodoByID)
	router.PUT("/todos/:id", updateTodo)
	router.DELETE("/todos/:id", deleteTodo)
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	fmt.Println("PORT: " + port)
	router.Run(":"+port)
}
