package rsaid

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// Simple library, not heavily focused on time - no need to pull in dependency.
// Using Tick variable for time.Now, which allows us to mock time.Now() in tests
// See: https://stackoverflow.com/questions/18970265/is-there-an-easy-way-to-stub-out-time-now-globally-during-test
var Tick = time.Now

type Gender int
type Citizenship int

const (
	GenderUnknown Gender = iota
	GenderMale
	GenderFemale
)

const (
	CitizenshipUnknown Citizenship = iota
	CitizenshipCitizen
	CitizenshipPermanentResident
	CitizenshipRefugee
)

type IdentityNumber struct {
	Value       string
	Citizenship Citizenship
	Gender      Gender
	DateOfBirth time.Time
}

// Parse parses a South African ID number string into details.
// Details include the ID number, citizenship, gender and date of birth.
// It returns the details and any errors encountered.
func Parse(number string) (IdentityNumber, error) {
	id := IdentityNumber{Value: number}
	if err := id.validate(); err != nil {
		return id, err
	}

	dob, err := id.dateOfBirth()
	if err != nil {
		return id, err
	}

	id.Value = number
	id.Citizenship = id.citizenship()
	id.Gender = id.gender()
	id.DateOfBirth = dob

	return id, nil
}

// validate determines if the given ID number is valid using the Luhn algorithm.
// It returns an error if the ID number is not a valid, has an invalid length or contains invalid characers or digits.
func (i IdentityNumber) validate() error {
	var sum int
	var alternate bool
	length := len(i.Value)
	if length != 13 {
		return errors.New("the provided south african id number does not equal 13 characters")
	}
	for j := length - 1; j > -1; j-- {
		mod, err := strconv.Atoi(string(i.Value[j]))
		if err != nil {
			return errors.New("the provided south african id number is not numeric")
		}
		// The 12th digit must always be 8
		if j == 11 && mod != 8 {
			return errors.New("the provided south african id number has an invalid and decommissioned race digit")
		}

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

// citizenship determines the citizenship status of the person in South Africa.
// Citizenship is calculated by using the 11th digit of the 13 digit ID number.
// Zero is considered a citizen, one a permanent resident, two a refugee, otherwise citizenship would be unknown.
// It returns the citizenship status of the person
func (i IdentityNumber) citizenship() Citizenship {
	switch cit, _ := strconv.Atoi(i.Value[10:11]); cit {
	case 0:
		return CitizenshipCitizen
	case 1:
		return CitizenshipPermanentResident
	case 2:
		return CitizenshipRefugee
	default:
		return CitizenshipUnknown
	}
}

// gender determines if the person is male or female.
// Gender is calculated by using the 7th digit of the 13 digit ID number.
// Zero to four is considered female, five to nine is considered male.
// It returns the gender of the person
func (i IdentityNumber) gender() Gender {
	// At this point, we can be assured that digit 7 is numeric
	gender, _ := strconv.Atoi(i.Value[6:7])
	if gender < 5 {
		return GenderFemale
	} else if gender < 9 {
		return GenderMale
	}
	return GenderUnknown
}

// dateOfBirth calculates the date of birth of the person.
// Date of birth is calculated by using the first 6 digits of the 13 digit ID number.
// The first pair of digits are the year, the second pair is the month and the third pair is the day.
// It returns the date of birth of the person and any errors encountered.
func (i IdentityNumber) dateOfBirth() (time.Time, error) {

	// Get current date along with assumed century
	curYear, curMonth, curDay := Tick().Date()
	curCentury := (curYear / 100) * 100

	// Get date values based off provided ID number
	// validate will have ensured we are working with numbers
	year, _ := strconv.Atoi(i.Value[0:2])
	year = curCentury + year
	month, _ := strconv.Atoi(i.Value[2:4])
	day, _ := strconv.Atoi(i.Value[4:6])

	// Only 16 years and above are eligible for an ID
	minYear := curYear - 16
	// Ensure the ID's DOB is not below 16 years from today, if so it's last century
	if year > minYear || (year == minYear && (month > int(curMonth) || month == int(curMonth) && day > curDay)) {
		year -= 100
	}

	loc, err := time.LoadLocation("Africa/Johannesburg")
	if err != nil {
		return Tick(), errors.New("could not load timezone to parse date of birth")
	}
	// Not using time.Date since it will still parse invalid dates. For example, 95/02/30 would parse to 95/03/02.
	dob, err := time.ParseInLocation("2006-01-02", fmt.Sprintf("%d-%02d-%02d", year, time.Month(month), day), loc)
	if err != nil {
		return Tick(), fmt.Errorf("cannot parse date of birth from id number: %s", err.Error())
	}

	return dob, nil
}
