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
