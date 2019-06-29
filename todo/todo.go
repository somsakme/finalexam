package todo

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/somsakme/finalexam/database"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type Todohandler struct{}

var todos = map[int]Todo{}

func (Todohandler) PostTodosHandler(c *gin.Context) {
	t := Todo{}
	fmt.Printf("befor post bind % #v\n", t)
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	fmt.Printf("After post bind % #v\n", t)

	db, err := database.GetDBConn()
	defer db.Close()

	query := `
		INSERT INTO customers (name,email,status) VALUES ($1,$2,$3) RETURNING id
		`
	var id int
	row := db.QueryRow(query, t.Name, t.Email, t.Status)
	err = row.Scan(&id)
	if err != nil {
		log.Fatal("Can't scan id", id)
	}
	fmt.Println("insert sucsess id:", id)
	t.ID = id
	c.JSON(http.StatusCreated, t)
}

func (Todohandler) GetTodosHandler(c *gin.Context) {
	db, err := database.GetDBConn()
	defer db.Close()

	stmt, err := db.Prepare("SELECT id,name,email,status FROM customers")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Prepare statment =": err.Error()})
		return
	}
	rows, _ := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Query =": err.Error()})
		return
	}

	todos := []Todo{}
	for rows.Next() {
		t := Todo{}
		err := rows.Scan(&t.ID, &t.Name, &t.Email, &t.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		todos = append(todos, t)
	}
	c.JSON(http.StatusOK, todos)
}

func (Todohandler) GetTodosByIdHandler(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDBConn()
	defer db.Close()

	stmt, _ := db.Prepare("SELECT id,name,email,status FROM customers WHERE id=$1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Prepare statment =": err.Error()})
		return
	}
	rows, _ := stmt.Query(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Query =": err.Error()})
		return
	}

	todos := []Todo{}
	t := Todo{}
	for rows.Next() {
		err := rows.Scan(&t.ID, &t.Name, &t.Email, &t.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		todos = append(todos, t)
	}
	c.JSON(http.StatusOK, t)
}

func (Todohandler) PutTodosByIDHanderler(c *gin.Context) {
	id := c.Param("id")
	t := Todo{}
	fmt.Printf("befor bind % #v\n", t)
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	fmt.Printf("After bind % #v\n", t)

	db, err := database.GetDBConn()
	defer db.Close()

	stmt, err := db.Prepare("UPDATE customers SET name=$2,email=$3,status=$4 WHERE id=$1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Prepare statment =": err.Error()})
		return
	}
	_, err = stmt.Query(id, t.Name, t.Email, t.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error update =": err.Error()})
		return
	}
	t.ID, err = strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Conv =": err.Error()})
		return
	}
	c.JSON(200, t)
}

func (Todohandler) DeleteTodosByIDHanderler(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDBConn()
	defer db.Close()

	stmt, _ := db.Prepare("DELETE FROM customers WHERE id=$1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Prepare statment =": err.Error()})
		return
	}
	_, err = stmt.Query(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Query =": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "customer deleted",
	})
}
