package api

import "net/http"

func RegisteredRoutes() {
	//associating routes with their handle functions(defined in /handlers)
	http.HandleFunc("/insert-data", InsertData)

	http.HandleFunc("/view-all-data", ViewAllData)

	http.HandleFunc("/get-rank", GetRank)

	http.HandleFunc("/update-score", UpdateScore)

	http.HandleFunc("/delete-record", DeleteRecord)
}
