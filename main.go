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

	//create database schema
	db.CreateSchema()

	//associating routes with their handle functions(defined in /handlers)
	http.HandleFunc("/insert-data", handlers.InsertData)

	http.HandleFunc("/view-all-data", handlers.ViewAllData)

	http.HandleFunc("/get-rank", handlers.GetRank)

	http.HandleFunc("/update-score", handlers.UpdateScore)

	http.HandleFunc("/delete-record", handlers.DeleteRecord)

	//Start an HTTP server on port 3000
	log.Fatal(http.ListenAndServe(":3000", nil))
}
