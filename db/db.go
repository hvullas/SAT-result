package db

import (
	"database/sql" // Import the "database/sql" package for database operations
	"fmt"          // Import the "fmt" package for string formatting
	"log"          // Import the "log" package for logging
	"os"           // Import the "os" package for accessing environment variables

	"github.com/joho/godotenv" // Import the "github.com/joho/godotenv" package for loading environment variables
	_ "github.com/lib/pq"      // Import the PostgreSQL driver for database/sql
)

// Declare a global variable for the database connection
var DB *sql.DB

// ConnectDB establishes a connection to the PostgreSQL database.
func ConnectDB() {

	// Load environment variables from a .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file couldn't load")
	}
	host := os.Getenv("HOST_NAME")
	port := os.Getenv("PORT")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DB_NAME")

	//constructs the connection string
	psqlinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//to open the database connection using the PostgreSQL driver
	DB, err = sql.Open("postgres", psqlinfo)
	if err != nil {
		log.Fatalln(err)
	}

	//logs a successful database connection
	log.Print("Connected to database")
}
