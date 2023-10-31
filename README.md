Project name : SAT-result

SAT-result is a Golang application that is built with PostgreSQL database and includes features as below
1. Insert data
2. View all data
3. Get rank
4. Update score
5. Delete one record

*Insert data handles the input data containg Name(unique identifier), address, city, country, pincode, sat_score of a condidate and stores in the database. This API includes validation for each input field.
  A constraint for sat score is implemented where sat score >30% is considerd as PASS.
*View all data responds with all the record data stored in the database.
*Get rank API gives the rank of the candidate with specified name. This dynamically changes if there is any update action on the sat_scores row of the database.
*Update score feature provide the functionality to update the SAT score of the candidate by name, and updates the pass_status column based on the passing score constraint.
*Delete record provides option to delete record of the candidate from the database using name.

Application is developed and tested using POSTMAN and Postman collections used are attached in the projects folder.

The project structure is as below:

||api
  -handler.go
  -routes.go
||db
  -db.go
||helperFuncs
  -helper.go
||models
  -models.go
|.env
|README.md
|go.mod
|go.sum
|main.go
|SATresults.postman_collection.json

*api package contains handlers.go file with all the handle functions and routes.go with insert data, get all data, get rank, update score and delete record routes.
*db package contains database connection functions using github.com/lib/pq driver and the database design schema for the application.
*helper.go file has helper functions and regex match patterns that are used for validation of the input.
*models.go has the structs defined in it which are used for handling input and output data.
*SATresults.postman_collection.json file contains the testing collection used while development.
