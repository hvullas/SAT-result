package handlers

import (
	"backend/db"
	"backend/models"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

func InsertData(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input models.SATresults
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	//validate name
	if len(input.Name) > 50 || len(input.Name) <= 0 {
		http.Error(w, "Invalid name", http.StatusBadRequest)
		return
	}

	//validate address
	if len(input.Address) > 100 || len(input.Address) <= 0 {
		http.Error(w, "Invalid address", http.StatusBadRequest)
		return
	}

	//validate city name
	if len(input.City) > 50 || len(input.City) <= 0 {
		http.Error(w, "Invalid city name", http.StatusBadRequest)
		return
	}

	//validate country name
	if len(input.Country) > 50 || len(input.Country) <= 0 {
		http.Error(w, "Invalid country name", http.StatusBadRequest)
		return
	}

	//validate pincode
	if len(input.Pincode) != 6 {
		http.Error(w, "Invalid postal code", http.StatusBadRequest)
		return
	}

	match, _ := regexp.MatchString(`\d{6}`, input.Pincode)
	if !match {
		http.Error(w, "Invalid postal code", http.StatusBadRequest)
		return
	}

	//validate SAT score
	if input.SATscore > 100.00 || input.SATscore < 0 {
		http.Error(w, "Invalid SAT score", http.StatusBadRequest)
		return
	}

	var passStatus bool
	if input.SATscore > 30.00 {
		passStatus = true
	}

	//insert to database
	insertQuery := `INSERT INTO results(name,address,city,country,pincode,sat_score,pass_status) values($1,$2,$3,$4,$5,$6,$7)`
	_, err = db.DB.Query(insertQuery, input.Name, input.Address, input.City, input.Country, input.Pincode, input.SATscore, passStatus)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)

}

func ViewAllData(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	getAllData := `SELECT * FROM results`

	var resultData []models.SATresults

	row, err := db.DB.Query(getAllData)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	for row.Next() {
		var data models.SATresults
		err = row.Scan(&data.Name, &data.Address, &data.City, &data.Country, &data.Pincode, &data.SATscore, &data.PassStatus)
		if err != nil {
			http.Error(w, "Scan error on data", http.StatusInternalServerError)
			return
		}
		resultData = append(resultData, data)
	}

	err = json.NewEncoder(w).Encode(resultData)
	if err != nil {
		http.Error(w, "Error encoding response body", http.StatusInternalServerError)
		return
	}

}

func GetRank(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input models.SATresults
	var rank models.Rank
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	if len(input.Name) == 0 || len(input.Name) > 50 {
		http.Error(w, "Invalid name", http.StatusBadRequest)
		return
	}

	rankQuery := `SELECT rank from (SELECT name,RANK() OVER ( ORDER BY sat_score DESC)rank FROM results) WHERE name=$1`

	err = db.DB.QueryRow(rankQuery, input.Name).Scan(&rank.Rank)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(rank)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

}

func UpdateScore(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input models.SATresults
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	if len(input.Name) == 0 || len(input.Name) > 50 {
		http.Error(w, "Invalid name", http.StatusBadRequest)
		return
	}

	if input.SATscore > 100.00 || input.SATscore < 0 {
		http.Error(w, "Invalid SAT score", http.StatusBadRequest)
		return
	}

	var passStatus bool
	if input.SATscore > 30.00 {
		passStatus = true
	}

	updateSATscoreQuery := `UPDATE results SET sat_score=$1,pass_status=$2 WHERE name=$3 RETURNING sat_score`
	err = db.DB.QueryRow(updateSATscoreQuery, input.SATscore, passStatus, input.Name).Scan(&input.SATscore)
	if err != nil {
		http.Error(w, "Error updating SAT score", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)

}

func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input models.SATresults
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	if len(input.Name) == 0 || len(input.Name) > 50 {
		http.Error(w, "Invalid name", http.StatusBadRequest)
		return
	}

	deleteRecordQuery := `DELETE FROM results WHERE name=$1`

	_, err = db.DB.Exec(deleteRecordQuery, input.Name)
	if err != nil {
		http.Error(w, "Error deleting the record", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)

}
