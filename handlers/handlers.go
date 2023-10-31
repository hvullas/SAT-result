package handlers

import (
	"backend/db"
	"backend/models"
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// InsertData is for handling reading input data from request and inserting it into the database
func InsertData(w http.ResponseWriter, r *http.Request) {

	//to log the called API
	log.Println(r.URL)

	//For returning JSON data in the response
	w.Header().Set("Content-Type", "application/json")

	//Ensures the request method must be POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input models.SATresults

	//decodes the request body into input variable
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	//validates for characters in the name are only alphabets and can have only space
	validNamePattern := regexp.MustCompile("^[A-Za-z\\s]+$")

	if !validNamePattern.MatchString(input.Name) {
		// Name is invalid
		http.Error(w, "Invalid name", http.StatusBadRequest)
		return
	}

	//len(input.Name)==0 ensure that the field is not empty
	//len(input.Name)>50 ensure that the name field does not contain unnecessary data
	if len(input.Name) > 50 || len(input.Name) <= 0 {
		http.Error(w, "Invalid name", http.StatusBadRequest)
		return
	}

	// capitalize the first letter of each word to ensure uniformity in the database.
	caser := cases.Title(language.English)
	input.Name = caser.String(input.Name)

	//validate input address
	if len(input.Address) > 100 || len(input.Address) <= 0 {
		http.Error(w, "Invalid address", http.StatusBadRequest)
		return
	}

	//validate input city name
	if len(input.City) > 50 || len(input.City) <= 0 {
		http.Error(w, "Invalid city name", http.StatusBadRequest)
		return
	}

	//validate input country name
	if len(input.Country) > 50 || len(input.Country) <= 0 {
		http.Error(w, "Invalid country name", http.StatusBadRequest)
		return
	}

	//validate input pincode
	if len(input.Pincode) != 6 {
		http.Error(w, "Invalid postal code", http.StatusBadRequest)
		return
	}

	//ensures that the pincode is a must be a 6digit number
	match, _ := regexp.MatchString(`\d{6}`, input.Pincode)
	if !match {
		http.Error(w, "Invalid postal code", http.StatusBadRequest)
		return
	}

	//validate input SAT score
	if input.SATscore > 100.00 || input.SATscore < 0 {
		http.Error(w, "Invalid SAT score", http.StatusBadRequest)
		return
	}

	//implements the constraint for pass or fail
	var passStatus bool

	//given constraint
	if input.SATscore > 30.00 {
		passStatus = true
	}

	//query for insertion of data into the database
	insertQuery := `INSERT INTO results(name,address,city,country,pincode,sat_score,pass_status) values($1,$2,$3,$4,$5,$6,$7)`

	//executing query to insert data into the table
	_, err = db.DB.Exec(insertQuery, input.Name, input.Address, input.City, input.Country, input.Pincode, input.SATscore, passStatus)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//To return success response with HTTP status code
	w.WriteHeader(200)

}

// viewAllData handle is for fetching all the data present in the results table in the database
func ViewAllData(w http.ResponseWriter, r *http.Request) {
	//to log the API call
	log.Println(r.URL)

	//for returning JSON data in response
	w.Header().Set("Content-Type", "application/json")

	//Ensures input request method is get
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//query to select all the data present in the results table in the database
	getAllData := `SELECT * FROM results`

	//to return array of record
	var resultData []models.SATresults

	//querying the database table for retriveing all the records present in the database
	row, err := db.DB.Query(getAllData)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	for row.Next() {

		//to read each record scanned by the query
		var data models.SATresults

		//reads into respective struct fields
		err = row.Scan(&data.Name, &data.Address, &data.City, &data.Country, &data.Pincode, &data.SATscore, &data.PassStatus)
		if err != nil {
			http.Error(w, "Scan error on data", http.StatusInternalServerError)
			return
		}

		//append each record to create an array of records to be returned with response
		resultData = append(resultData, data)
	}

	//Encodeing  resultData array in json
	err = json.NewEncoder(w).Encode(resultData)
	if err != nil {
		http.Error(w, "Error encoding response body", http.StatusInternalServerError)
		return
	}

}

// GetRank handle for responding with the rank obtained by the candidate
func GetRank(w http.ResponseWriter, r *http.Request) {

	//to log the api call
	log.Println(r.URL)

	//for returning response in JSON format
	w.Header().Set("Content-Type", "application/json")

	//To ensure input request method is get
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//to read the input name
	//struct omits all other fields of the SAtresults struct except name
	var input models.SATresults

	//to send rank as the response
	var rank models.Rank

	//for decoding the request body
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	//len(input.Name)==0 ensure that the field is not empty
	//len(input.Name)>50 ensure that the name field does not contain unnecessary data
	if len(input.Name) == 0 || len(input.Name) > 50 {
		http.Error(w, "Invalid name", http.StatusBadRequest)
		return
	}

	//sql query string for selecting rank from the records present in the results table in database
	rankQuery := `SELECT rank from (SELECT name,RANK() OVER ( ORDER BY sat_score DESC)rank FROM results) WHERE name=$1`

	//querying the database table for retriveing rank of the student
	err = db.DB.QueryRow(rankQuery, input.Name).Scan(&rank.Rank)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Encode the response in JSON
	err = json.NewEncoder(w).Encode(rank)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

}

// handle function to update the SAT score of the candidate
func UpdateScore(w http.ResponseWriter, r *http.Request) {
	//to log the API call
	log.Println(r.URL)

	//for ensuring that the response in in JSON format
	w.Header().Set("Content-Type", "application/json")

	//Ensure that the input request method is PUT
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//for reading the inputs name and sat_score
	var input models.SATresults

	//To decode the request body
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	//len(input.Name)==0 ensure that the field is not empty
	//len(input.Name)>50 ensure that the name field does not contain unnecessary data
	if len(input.Name) == 0 || len(input.Name) > 50 {
		http.Error(w, "Invalid name", http.StatusBadRequest)
		return
	}

	//ensuring the input SAT score in valid and stays between 0.00 and 100.00
	if input.SATscore > 100.00 || input.SATscore < 0 {
		http.Error(w, "Invalid SAT score", http.StatusBadRequest)
		return
	}

	//for checking and determining pass or fail based on the SAT score
	var passStatus bool
	if input.SATscore > 30.00 {
		passStatus = true
	}

	//sql query to update the sat_score of a candidate
	updateSATscoreQuery := `UPDATE results SET sat_score=$1,pass_status=$2 WHERE name=$3 RETURNING sat_score`

	//querying the database for updating the records information
	err = db.DB.QueryRow(updateSATscoreQuery, input.SATscore, passStatus, input.Name).Scan(&input.SATscore)
	if err != nil {
		http.Error(w, "Error updating SAT score", http.StatusInternalServerError)
		return
	}

	//return with the success response
	w.WriteHeader(200)

}

// For deleting the record i the database table with the given name of the candidate
func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	//to log the API call
	log.Println(r.URL)

	//To ensure that the response is in JSON
	w.Header().Set("Content-Type", "application/json")

	//To ensure request methos is DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//to read the name from the input
	var input models.SATresults

	//decode the JSON request body into variable
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	//len(input.Name)==0 ensure that the field is not empty
	//len(input.Name)>50 ensure that the name field does not contain unnecessary data
	if len(input.Name) == 0 || len(input.Name) > 50 {
		http.Error(w, "Invalid name", http.StatusBadRequest)
		return
	}

	//sql query to delete record from the results table with specified name
	deleteRecordQuery := `DELETE FROM results WHERE name=$1`

	//Query the database and delete the record
	_, err = db.DB.Exec(deleteRecordQuery, input.Name)
	if err != nil {
		http.Error(w, "Error deleting the record", http.StatusInternalServerError)
		return
	}

	//returns with success response
	w.WriteHeader(200)

}
