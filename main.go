package main

import (
	"backend/api" //package with registered routes
	"backend/db"  //package for database operations
	"log"         //Imported for logging
	"net/http"    //package for creating an http server
)

func main() {

	//establishes connection with database and ensures it closes when program exits
	db.ConnectDB()
	defer db.DB.Close()

	//create database schema
	db.CreateDatabaseSchema()

	//for accessing registered routes
	api.RegisteredRoutes()

	//Start an HTTP server on port 3000
	log.Fatal(http.ListenAndServe(":3000", nil))
}
