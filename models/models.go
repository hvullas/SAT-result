package models

type SATresults struct {
	Name       string  `json:"name"`
	Address    string  `json:"address"`
	City       string  `json:"city"`
	Country    string  `json:"country"`
	Pincode    string  `json:"pincode"`
	SATscore   float64 `json:"sat_score"`
	PassStatus bool    `json:"pass_status,omitempty"`
}
