package rsaid

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Gender int

const (
	GenderUnknown Gender = iota
	GenderMale
	GenderFemale
)

type IdentityNumber struct {
	raw  string
}

func Parse(idNumber string) (IdentityNumber, error) {
	idNum := IdentityNumber{raw: idNumber}
	if err := idNum.validate(); err != nil {
		return idNum, err
	}

	return idNum, nil
}

func (i IdentityNumber) validate() error {
	var sum int
	var alternate bool
	length := len(i.raw)
	if length != 13 {
		return errors.New("the provided south african id number does not equal 13 characters")
	}
	for j := length - 1; j > -1; j-- {
		mod, err := strconv.Atoi(string(i.raw[j]))
		if err != nil {
			return errors.New("the provided south african id number is not numeric")
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

func (i IdentityNumber) IsCitizen() bool {
	return i.raw[10:11] == "0"
}

func (i IdentityNumber) Gender() Gender {
	// At this point, we can be assured that digit 7 is numeric
	gender, _ := strconv.Atoi(i.raw[6:7])
	if gender < 5 {
		return GenderFemale
	} else if gender < 9 {
		return GenderMale
	}
	return GenderUnknown
}

func (i IdentityNumber) DateOfBirth() (time.Time, error) {
	// Get current date along with assumed century
	CurrentYear, CurrentMonth, CurrentDay := time.Now().Date()
	CurrentCentury := (CurrentYear / 100) * 100

	// Get date values based off provided ID number
	// validate will have ensured we are working with numbers
	ProvidedYear, _ := strconv.Atoi(i.raw[0:2])
	ProvidedYear = CurrentCentury + ProvidedYear
	ProvidedMonth, _ := strconv.Atoi(i.raw[2:4])
	ProvidedDay, _ := strconv.Atoi(i.raw[4:6])

	// Only 16 years and above are eligible for an ID
	EligibleYear := CurrentYear - 16
	// Ensure the ID's DOB is not below 16 years from today, if so it's last century
	if ProvidedYear > EligibleYear || (ProvidedYear == EligibleYear && (ProvidedMonth > int(CurrentMonth) || ProvidedMonth == int(CurrentMonth) && ProvidedDay > CurrentDay)) {
		ProvidedYear -= 100
	}

	loc, err := time.LoadLocation("Africa/Johannesburg")
	if err != nil {
		return time.Now(), errors.New("could not load timezone to parse date of birth")
	}
	// Not using time.Date since it will still parse invalid dates. For example, 95/02/30 would parse to 95/03/02.
	dob, err := time.ParseInLocation("2006-01-02", fmt.Sprintf("%d-%02d-%02d", ProvidedYear, time.Month(ProvidedMonth), ProvidedDay), loc)
	if err != nil {
		return time.Now(), fmt.Errorf("cannot parse date of birth from id number: %s", err.Error())
	}

	return dob, nil
}
