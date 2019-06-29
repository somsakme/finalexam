package main

import (
	"fmt"
	"net/http"

	"github.com/somsakme/finalexam/database"
	"github.com/somsakme/finalexam/todo"

	"github.com/gin-gonic/gin"
)

func authMiddlware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "token2019" {
		fmt.Println(http.StatusUnauthorized)
		c.JSON(http.StatusUnauthorized, gin.H{"error Token": http.StatusText(http.StatusUnauthorized)})
		c.Abort()
		return
	}

	c.Next()
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(authMiddlware)

	s := todo.Todohandler{}
	r.POST("/customers", s.PostTodosHandler)
	r.GET("/customers", s.GetTodosHandler)
	r.GET("/customers/:id", s.GetTodosByIdHandler)
	r.PUT("/customers/:id", s.PutTodosByIDHanderler)
	r.DELETE("/customers/:id", s.DeleteTodosByIDHanderler)

	return r
}

func main() {
	database.CreateDatabase()

	r := setupRouter()
	r.Run(":2019")
}
