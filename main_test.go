package rsaid_test

import (
	"errors"
	"testing"

	"github.com/jacovdloo/rsaid"
)

func Test_Parse(t *testing.T) {

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
			name:    "Invalid DOB",
			id:      "9502305120004",
			wantErr: errors.New("cannot parse date of birth from id number: parsing time \"1995-02-30\": day out of range"),
		},
		{
			name: "Valid ID",
			id:   "9506245120008",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := rsaid.Parse(tt.id)
			if tt.wantErr != nil && err != nil && tt.wantErr.Error() != err.Error() {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_Gender(t *testing.T) {

	tests := []struct {
		name string
		id   string
		want rsaid.Gender
	}{
		{
			name: "Is Male",
			id:   "9506245120008",
			want: rsaid.GenderMale,
		},
		{
			name: "Is Female",
			id:   "9506244120009",
			want: rsaid.GenderFemale,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idNum, _ := rsaid.Parse(tt.id)
			if tt.want != idNum.Gender(){
				t.Errorf("Does not determine gender correctly")
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
			idNum, _ := rsaid.Parse(tt.id)
			if tt.want != idNum.IsCitizen() {
				t.Errorf("Does not determine citizenship correctly")
			}
		})
	}
}

func Test_DateOfBirth(t *testing.T) {

	// These tests will fail once time.Now() causes the age to not be valid.  Best to mock time.
	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "Between 16 and 100",
			id:   "9506245120008",
			want: "1995-06-24 00:00:00 +0200 SAST",
		},
		{
			name: "Over 100",
			id:   "2201014800082",
			want: "1922-01-01 00:00:00 +0200 SAST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idNum, _ := rsaid.Parse(tt.id)
			dob := idNum.DateOfBirth()
			if tt.want != dob.String() {
				t.Errorf("Does not determine date of birth correctly. dob = %v", dob)
			}
		})
	}
}
