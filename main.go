package main

import (
	"backend/db"       //package for database operations
	"backend/handlers" //package for request handlers
	"log"              //Imported for logging
	"net/http"         //package for creating an http server
)

func main() {

	//establishes connection with database and ensures it closes when program exits
	db.ConnectDB()
	defer db.DB.Close()

	//SQL query that creates table 'results' if it doest not exist in the database
	createResultsTable := ` CREATE TABLE IF NOT EXISTS results(
		name VARCHAR(50) PRIMARY KEY,
		address VARCHAR(100) NOT NULL,
		city VARCHAR(50) NOT NULL,
		country VARCHAR(50) NOT NULL,
		pincode VARCHAR(10) NOT NULL,
		sat_score NUMERIC(5,2) NOT NULL,
		pass_status BOOLEAN
	)`

	_, err := db.DB.Exec(createResultsTable)
	if err != nil {
		log.Panicln("Error creating table in the database")
		return
	}

	//associating routes with their handle functions(defined in /handlers)
	http.HandleFunc("/insert-data", handlers.InsertData)

	http.HandleFunc("/view-all-data", handlers.ViewAllData)

	http.HandleFunc("/get-rank", handlers.GetRank)

	http.HandleFunc("/update-score", handlers.UpdateScore)

	http.HandleFunc("/delete-record", handlers.DeleteRecord)

	//Start an HTTP server on port 3000
	log.Fatal(http.ListenAndServe(":3000", nil))
}
