package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func CreateDatabase() {
	db, err := GetDBConn()
	defer db.Close()

	createTb := `
	CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	);
	`
	_, err = db.Exec(createTb)
	if err != nil {
		log.Fatal("Canot create database", err.Error())
	}
	//	fmt.Println("Create Database completed")
}

func GetDBConn() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	return db, nil
}
