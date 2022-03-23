package rsaid_test

import (
	"errors"
	"testing"

	"github.com/jacovdloo/rsaid"
)

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
			wantErr: errors.New("the provided south african id number is not numeric"),
		},
		{
			name:    "Random Letters",
			id:      "randomletters",
			wantErr: errors.New("the provided south african id number is not numeric"),
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
			id:   "9506245120008",
			want: int(rsaid.GenderMale),
		},
		{
			name: "Is Female",
			id:   "9506244120009",
			want: int(rsaid.GenderFemale),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gender, err := rsaid.Gender(tt.id)
			if int(gender) != tt.want || err != nil {
				if err != nil {
					t.Errorf("Does not determine gender correctly. error = %v", err)
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
			id:   "9506245120008",
			want: true,
		},
		{
			name: "Permanent Resident",
			id:   "9506245120107",
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

	maleCitizen := "9506245120008"

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

	// Invalid DOB - 95/02/30
	if _, err = rsaid.DateOfBirth("9502305120004"); err == nil {
		t.Errorf("Does not determine date of birth correctly")
	}
}

func Test_Parse(t *testing.T) {

	maleCitizen := "9506245120008"

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
