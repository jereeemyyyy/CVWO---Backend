package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var DB *sql.DB

// InitDB initializes the database connection.
func InitDB() {
	host := os.Getenv("SQL_HOST")
	username := os.Getenv("DBUSER")
	password := os.Getenv("DBPASS")
	album := os.Getenv("ALBUM")
	connStr := fmt.Sprintf("host=%s user=%s dbname=%s password=%s", host, username, album, password)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	pingErr := DB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!")
}

