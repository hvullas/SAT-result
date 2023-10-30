package main

import (
	"backend/db"
	"backend/handlers"
	"log"
	"net/http"
)

func main() {
	db.ConnectDB()
	defer db.DB.Close()

	createTable := ` CREATE TABLE IF NOT EXISTS results(
		name VARCHAR(50) PRIMARY KEY,
		address VARCHAR(100) NOT NULL,
		city VARCHAR(50) NOT NULL,
		country VARCHAR(50) NOT NULL,
		pincode VARCHAR(10) NOT NULL,
		sat_score NUMERIC(5,2) NOT NULL,
		pass_status BOOLEAN
	)`

	_, err := db.DB.Exec(createTable)
	if err != nil {
		log.Panicln("Error creating table in the database")
		return
	}

	http.HandleFunc("/insert-data", handlers.InsertData)

	http.HandleFunc("/view-all-data", handlers.ViewAllData)

	http.HandleFunc("/get-rank", handlers.GetRank)

	http.HandleFunc("/update-score", handlers.UpdateScore)

	http.HandleFunc("/delete-record", handlers.DeleteRecord)

	log.Fatal(http.ListenAndServe(":3000", nil))
}
