package rsaid_test

import (
	"errors"
	"testing"

	"github.com/jacovdloo/rsaid"
)

var nonCitizen = "9506245120107"
var maleCitizen = "9506245120008"
var femaleCitizen = "9506244120009"
var invalidDateOfBirth = "9502305120004"

func Test_IsValid(t *testing.T) {

	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{
			name:    "Too Short",
			id:      "950624",
			wantErr: errors.New("the provided south african id number does not equal 13 characters"),
		},
		{
			name:    "Too Long",
			id:      "95062451200010",
			wantErr: errors.New("the provided south african id number does not equal 13 characters"),
		},
		{
			name:    "Seven Letters",
			id:      "letters",
			wantErr: errors.New("the provided south african id number does not equal 13 characters"),
		},
		{
			name:    "Invalid Character",
			id:      "950624G120008",
			wantErr: errors.New("the provided south african id number should be numeric"),
		},
		{
			name:    "Random Letters",
			id:      "randomletters",
			wantErr: errors.New("the provided south african id number should be numeric"),
		},
		{
			name:    "Invalid ID",
			id:      "9506245120009",
			wantErr: errors.New("the provided south african id number is not valid"),
		},
		{
			name: "Valid ID",
			id:   "9506245120008",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rsaid.IsValid(tt.id)
			if tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("IsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_Gender(t *testing.T) {

	tests := []struct {
		name string
		id   string
		want int
	}{
		{
			name: "Is Male",
			id:   maleCitizen,
			want: int(rsaid.GenderMale),
		},
		{
			name: "Is Female",
			id:   femaleCitizen,
			want: int(rsaid.GenderFemale),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gender, err := rsaid.Gender(tt.id)
			if int(gender) != tt.want || err != nil {
				if err != nil {
					t.Errorf("Does not determine gender correctly: %v", err)
				} else {
					t.Errorf("Does not determine gender correctly")
				}
			}
		})
	}
}

func Test_IsCitizen(t *testing.T) {

	tests := []struct {
		name string
		id   string
		want bool
	}{
		{
			name: "Citizenship",
			id:   maleCitizen,
			want: true,
		},
		{
			name: "Permanent Resident",
			id:   nonCitizen,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cit, err := rsaid.IsCitizen(tt.id)
			if cit != tt.want || err != nil {
				if err != nil {
					t.Errorf("Does not determine citizenship correctly. error = %v", err)
				} else {
					t.Errorf("Does not determine citizenship correctly")
				}
			}
		})
	}
}

func Test_DateOfBirth(t *testing.T) {

	dob, err := rsaid.DateOfBirth(maleCitizen)

	if err != nil {
		t.Errorf("Does not determine date of birth correctly. error = %v", err)
		return
	}

	if dob.Year() != 1995 {
		t.Errorf("Does not determine date of birth year correctly")
	}
	if dob.Month() != 6 {
		t.Errorf("Does not determine date of birth month correctly")
	}
	if dob.Day() != 24 {
		t.Errorf("Does not determine date of birth day correctly")
	}

	_, err = rsaid.DateOfBirth(invalidDateOfBirth)
	if err == nil {
		t.Errorf("Does not determine date of birth correctly")
	}
}

func Test_Parse(t *testing.T) {

	person, err := rsaid.Parse(maleCitizen)

	if err != nil {
		t.Errorf("Does not parse ID number correctly: %s", err.Error())
	}
	if person.Gender != rsaid.GenderMale {
		t.Errorf("Does not parse gender correctly")
	}
	if person.Citizen != true {
		t.Errorf("Does not parse citizenship correctly")
	}
	if person.DOB.Year() != 1995 || person.DOB.Month() != 6 || person.DOB.Day() != 24 {
		t.Errorf("Does not parse birth date correctly")
	}
}
