package rsaid

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type IDGender string

const (
	Unknown IDGender = "unknown"
	Male    IDGender = "male"
	Female  IDGender = "female"
)

type Details struct {
	DOB     time.Time
	Gender  IDGender
	Citizen bool
}

// IsValid determines if the given ID number is valid using the Luhn algorithm.
func IsValid(id string) error {
	var sum int
	var alternate bool
	length := len(id)
	if length != 13 {
		return errors.New("the provided south african id number does not equal 13 characters")
	}
	for i := length - 1; i > -1; i-- {
		mod, _ := strconv.Atoi(string(id[i]))
		if alternate {
			mod *= 2
			if mod > 9 {
				mod = (mod % 10) + 1
			}
		}
		alternate = !alternate
		sum += mod
	}
	if sum%10 == 0 {
		return nil
	} else {
		return errors.New("the provided south african id number is not valid")
	}
}

// Gender determines if the person is male or female.
// Gender is calculated by using the 7th digit in the 13 digit ID number.
// Zero to four is considered female, five to nine is considered male.
// It returns either "male" or "female" and any errors encountered.
func Gender(id string) (IDGender, error) {
	if err := IsValid(id); err != nil {
		return Unknown, err
	}
	gender, _ := strconv.Atoi(id[6:7])
	if gender < 5 {
		return Female, nil
	}
	return Male, nil
}

// IsCitizen determines if the person is a South African citizen.
// Citizenship is calculated by using the 11th digit in the 13 digit ID number.
// Zero is considered a citizen, otherwise, it is considered a permanent resident.
// It returns true if the person is a citizen and any errors encountered.
func IsCitizen(id string) (bool, error) {
	if err := IsValid(id); err != nil {
		return false, err
	}
	citizenCode := id[10:11]
	return citizenCode == "0", nil
}

// DateOfBirth calculates the date of birth of the person.
// Date of birth is calculated by using the first 6 digits in the 13 digit ID number.
// The first pair of digits are the year, the second pair is the month and the third pair is the day.
// It returns the date of birth and any errors encountered.
func DateOfBirth(id string) (time.Time, error) {
	if err := IsValid(id); err != nil {
		return time.Now(), err
	}
	// Get current date along with assumed century
	CurrentYear, CurrentMonth, CurrentDay := time.Now().Date()
	CurrentCentury := (CurrentYear / 100) * 100

	// Get date values based off provided ID number
	ProvidedYear, _ := strconv.Atoi(id[0:2])
	ProvidedYear = CurrentCentury + ProvidedYear
	ProvidedMonth, _ := strconv.Atoi(id[2:4])
	ProvidedDay, _ := strconv.Atoi(id[4:6])

	// Only 16 years and above are eligible for an ID
	EligibleYear := CurrentYear - 16
	// Ensure the ID's DOB is not below 16 years from today, if so it's last century
	if ProvidedYear > EligibleYear || (ProvidedYear == EligibleYear && (ProvidedMonth > int(CurrentMonth) || ProvidedMonth == int(CurrentMonth) && ProvidedDay > CurrentDay)) {
		ProvidedYear -= 100
	}
	// See https://pkg.go.dev/time for standard templates
	// Use South African Standard Time (UTC+2).
	// Registration of ID would have happened in South Africa, regardless of citizenship or permanent resident.
	loc, _ := time.LoadLocation("Africa/Johannesburg")
	dob, err := time.ParseInLocation("2006-01-02", fmt.Sprintf("%d-%02d-%02d", ProvidedYear, time.Month(ProvidedMonth), ProvidedDay), loc)
	if err != nil {
		return time.Now(), errors.New("cannot parse date of birth from id number")
	}
	return dob, nil
}

// Parse derives details from a South African ID number.
// Details include gender, citizenship, and date of birth.
// It returns the details and any errors encountered.
func Parse(id string) (Details, error) {
	if err := IsValid(id); err != nil {
		return Details{}, err
	}
	dob, _ := DateOfBirth(id)
	gender, _ := Gender(id)
	citizen, _ := IsCitizen(id)
	return Details{
		DOB:     dob,
		Gender:  gender,
		Citizen: citizen,
	}, nil
}
