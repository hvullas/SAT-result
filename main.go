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

	http.HandleFunc("/insert-data", handlers.InsertData)

	http.HandleFunc("/view-all-data", handlers.ViewAllData)

	http.HandleFunc("/get-rank", handlers.GetRank)

	http.HandleFunc("/update-score", handlers.UpdateScore)

	http.HandleFunc("/delete-record", handlers.DeleteRecord)

	log.Fatal(http.ListenAndServe(":3000", nil))
}
