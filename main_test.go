package rsaid

import (
	"errors"
	"testing"
	"time"
)

func TestParse(t *testing.T) {

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
			id:      "95062451201802",
			wantErr: errors.New("the provided south african id number does not equal 13 characters"),
		},
		{
			name:    "Seven Letters",
			id:      "letters",
			wantErr: errors.New("the provided south african id number does not equal 13 characters"),
		},
		{
			name:    "Invalid Character",
			id:      "950624G120081",
			wantErr: errors.New("the provided south african id number is not numeric"),
		},
		{
			name:    "Random Letters",
			id:      "randomletters",
			wantErr: errors.New("the provided south african id number is not numeric"),
		},
		{
			name:    "Invalid ID",
			id:      "9506245120082",
			wantErr: errors.New("the provided south african id number is not valid"),
		},
		{
			name:    "Invalid DOB",
			id:      "9502305120087",
			wantErr: errors.New(`cannot parse date of birth from id number: parsing time "1995-02-30": day out of range`),
		},
		{
			name:    "Invalid Race",
			id:      "9506245120001",
			wantErr: errors.New("the provided south african id number has an invalid and decommissioned race digit"),
		},
		{
			name: "Valid ID",
			id:   "9506245120081",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.id)
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
		want Gender
	}{
		{
			name: "Is Male",
			id:   "9506245120081",
			want: GenderMale,
		},
		{
			name: "Is Female",
			id:   "9506244120082",
			want: GenderFemale,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, _ := Parse(tt.id)
			gen := id.Gender()
			if tt.want != gen {
				t.Errorf("Does not determine gender correctly. want gender = %v, got gender = %v", tt.want, gen)
			}
		})
	}
}

func Test_Citizen(t *testing.T) {

	tests := []struct {
		name string
		id   string
		want Citizenship
	}{
		{
			name: "Citizen",
			id:   "9506245120081",
			want: CitizenshipCitizen,
		},
		{
			name: "Permanent Resident",
			id:   "9506245120180",
			want: CitizenshipResident,
		},
		{
			name: "Refugee",
			id:   "9901019999283",
			want: CitizenshipRefugee,
		},
		{
			name: "Unknown Residency",
			id:   "9901019999382",
			want: CitizenshipUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, _ := Parse(tt.id)
			cit := id.Citizenship()
			if tt.want != cit {
				t.Errorf("Does not determine citizenship correctly. want citizenship = %v, got citizenship = %v", tt.want, cit)
			}
		})
	}
}

func Test_DateOfBirth(t *testing.T) {

	// Mock the clock
	Tick = func() time.Time { return time.Date(2020, time.Month(6), 20, 0, 0, 0, 0, time.UTC) }

	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "Between 16 and 100",
			id:   "9506245120081",
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
			id, _ := Parse(tt.id)
			dob := id.DateOfBirth()
			if tt.want != dob.String() {
				t.Errorf("Does not determine date of birth correctly. want dob = %v, got dob = %v", tt.want, dob)
			}
		})
	}
}
