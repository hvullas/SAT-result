package models

//This struct is to read and output the information,
// omitempty helps reuse of the struct and handling unnecessary encoding and decoding of fields
type SATresults struct {
	Name       string  `json:"name,omitempty"`
	Address    string  `json:"address,omitempty"`
	City       string  `json:"city,omitempty"`
	Country    string  `json:"country,omitempty"`
	Pincode    string  `json:"pincode,omitempty"`
	SATscore   float64 `json:"sat_score,omitempty"`
	PassStatus bool    `json:"pass_status,omitempty"`
}

//To output the rank of the specified name
type Rank struct {
	Rank int64 `json:"rank"`
}
