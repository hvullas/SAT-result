package api

import (
	"backend/db"
	helper "backend/helperFuncs"
	"backend/models"
	"database/sql"
	"encoding/json"
	"net/http"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// custom errors
const (
	ErrInvalidName         = "Invalid name"
	ErrInvalidAddress      = "Invalid address"
	ErrInvalidCity         = "Invalid city name"
	ErrInvalidCountry      = "Invalid country name"
	ErrInvalidPincode      = "Invalid postal code"
	ErrInvalidSATScore     = "Invalid SAT score"
	ErrDatabaseError       = "Internal server error"
	ErrMethodNotAllowed    = "Method not allowed"
	ErrDecodingRequestBody = "Error decoding JSON request body"
	ErrEncodingResponse    = "Error encoding response"
)

// InsertData is for handling reading input data from request and inserting it into the database
func InsertData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var input models.SATresults

	//decodes the request body into input variable
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	if !helper.ValidName(input.Name) {
		http.Error(w, ErrInvalidName, http.StatusBadRequest)
		return
	}
	// capitalize the first letter of each word to ensure uniformity in the database.
	caser := cases.Title(language.English)
	input.Name = caser.String(input.Name)

	if !helper.ValidAddress(input.Address) {
		http.Error(w, ErrInvalidAddress, http.StatusBadRequest)
		return
	}

	if !helper.ValidCityName(input.City) {
		http.Error(w, ErrInvalidCity, http.StatusBadRequest)
		return
	}

	if !helper.ValidCountryName(input.Country) {
		http.Error(w, ErrInvalidCountry, http.StatusBadRequest)
		return
	}

	if !helper.ValidSATscore(input.SATscore) {
		http.Error(w, ErrInvalidSATScore, http.StatusBadRequest)
		return
	}

	//candidate passes if score is more than 30%
	var passStatus bool
	if input.SATscore > 30.00 {
		passStatus = true
	}

	insertQuery := `INSERT INTO results(name,address,city,country,pincode,sat_score,pass_status) values($1,$2,$3,$4,$5,$6,$7)`
	_, err = db.DB.Exec(insertQuery, input.Name, input.Address, input.City, input.Country, input.Pincode, input.SATscore, passStatus)
	if err != nil {
		http.Error(w, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)

}

// viewAllData handle is for fetching all the data present in the results table in the database
// response is an array of records of candidates
func ViewAllData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	getAllData := `SELECT * FROM results`
	var resultData []models.SATresults
	row, err := db.DB.Query(getAllData)
	if err != nil {
		http.Error(w, ErrDatabaseError, http.StatusInternalServerError)
		return
	}
	for row.Next() {
		var data models.SATresults
		err = row.Scan(&data.Name, &data.Address, &data.City, &data.Country, &data.Pincode, &data.SATscore, &data.PassStatus)
		if err != nil {
			http.Error(w, ErrDatabaseError, http.StatusInternalServerError)
			return
		}
		resultData = append(resultData, data)
	}

	err = json.NewEncoder(w).Encode(resultData)
	if err != nil {
		http.Error(w, ErrEncodingResponse, http.StatusInternalServerError)
		return
	}

}

// GetRank handle for responding with the rank obtained by the candidate
func GetRank(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}
	var input models.SATresults
	//to send rank as the response
	var rank models.Rank
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	if !helper.ValidName(input.Name) {
		http.Error(w, ErrInvalidName, http.StatusBadRequest)
		return
	}

	rankQuery := `SELECT rank from (SELECT name,RANK() OVER ( ORDER BY sat_score DESC)rank FROM results) WHERE name=$1`
	err = db.DB.QueryRow(rankQuery, input.Name).Scan(&rank.Rank)
	if err != nil {
		if err == sql.ErrNoRows { //Handle error for name not found in database
			http.Error(w, ErrInvalidName, http.StatusBadRequest)
			return
		}
		http.Error(w, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(rank)
	if err != nil {
		http.Error(w, ErrEncodingResponse, http.StatusInternalServerError)
		return
	}

}

// handle function to update the SAT score of the candidate
func UpdateScore(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPut {
		http.Error(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var input models.SATresults
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	if !helper.ValidName(input.Name) {
		http.Error(w, ErrInvalidName, http.StatusBadRequest)
		return
	}

	if !helper.ValidSATscore(input.SATscore) {
		http.Error(w, ErrInvalidSATScore, http.StatusBadRequest)
		return
	}

	//for checking and determining pass or fail based on the SAT score
	var passStatus bool
	if input.SATscore > 30.00 {
		passStatus = true
	}

	updateSATscoreQuery := `UPDATE results SET sat_score=$1,pass_status=$2 WHERE name=$3 RETURNING sat_score`
	err = db.DB.QueryRow(updateSATscoreQuery, input.SATscore, passStatus, input.Name).Scan(&input.SATscore)
	if err != nil {
		http.Error(w, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)

}

// For deleting the record i the database table with the given name of the candidate
func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		http.Error(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}
	var input models.SATresults
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	if !helper.ValidName(input.Name) {
		http.Error(w, ErrInvalidName, http.StatusBadRequest)
		return
	}

	deleteRecordQuery := `DELETE FROM results WHERE name=$1`
	_, err = db.DB.Exec(deleteRecordQuery, input.Name)
	if err != nil {
		http.Error(w, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)

}
