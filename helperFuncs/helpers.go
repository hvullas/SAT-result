package helper

import (
	"regexp"
)

func MatchPattern(structfield string) bool {
	//validates for characters in the name are only alphabets and can have only space
	validNamePattern := regexp.MustCompile("^[A-Za-z\\s]+$")

	return validNamePattern.MatchString(structfield)
}

// names of candidate,country name, city name is considered to be less than 50 characters and contains only alphabet and space
func ValidName(name string) bool {
	return len(name) != 0 && len(name) < 50 && MatchPattern(name)
}

func ValidCountryName(country string) bool {
	return len(country) < 50 && len(country) != 0 && MatchPattern(country)
}

func ValidCityName(city string) bool {
	return len(city) < 50 && len(city) != 0 && MatchPattern(city)
}

// considered to have less than 100 characters
func ValidAddress(address string) bool {
	return len(address) < 100 && len(address) != 0
}

// pincode is considered to have only 6 digits number
func ValidPincode(pincode string) bool {
	if len(pincode) != 6 {
		return false
	}

	//ensures that the pincode is a must be a 6digit number
	match, _ := regexp.MatchString(`\d{6}`, pincode)
	return match
}

// SAT score is read from input as percentage from 0.00 to 100.00
func ValidSATscore(score float64) bool {
	return score <= 100.00 || score >= 0.00
}
