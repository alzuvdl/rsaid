# 🇿🇦 South African ID Validator for Go

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/mod/github.com/jacovdloo/rsaid)

Validate ID numbers for the Republic of South Africa - taking eligibility age (16 years) into account.

The following information can be derived if the ID Number is valid:

- Gender
- Citizenship
- Date of birth

## Details

A South African ID number is a 13-digit number which is defined by the following format: `YYMMDDSSSSCAZ`.

The first 6 digits `YYMMDD` are based on the person's date of birth. For example, 24 June 1995 is displayed as `950624`

`{950624} SNNNCAZ`

> Although rare, it can happen that someone’s birth date does not correspond with their ID number.

The next digit `S` is used to define the person's gender. Females are assigned numbers in the range `0-4` and males from `5-9`.

`{950624} {5} NNNCAZ`

The next three digits `NNN` indicates the person's position of birth in the registry based on the date of birth and gender.

If for example, the number is 120, it means that person was the 120th person to be registered as having been born on that particular day, for that gender.

`{950624} {5} {120} CAZ`

The next digit `C` indicates the person's citizenship. `0` denoting that the person was born a South African citizen, `1` denoting that the person is a permanent resident and `2` denoting that the person is a refugee.

`{950624} {5} {120} {0} AZ`

The next digit `A` was used until the late 1980s to indicate a person’s race. This has been eliminated and old ID numbers were reissued to remove this.

> Pre-democracy classifications.

> Before the race group classification was abandoned, this is what digit `A` in the ID number indicated:

- 0 — White
- 1 – Cape Coloured
- 2 – Malay
- 3 – Griqua
- 4 – Chinese
- 5 – Indian
- 6 – Other Asian
- 7 – Other Coloured

All new ID numbers must have an `A` digit value of `8`

`{950624} {5} {120} {0} {8} Z`

The last digit `Z` is a checksum digit, used to check that the number sequence is accurate using the [Luhn](https://en.wikipedia.org/wiki/Luhn_algorithm) algorithm.

`{950624} {5} {120} {0} {8} {1}`

So, ID number `9506245120081` will parse to the `120th` `male` South African `citizen` born/registered on the `24th of June, 1995`

## Install

```go
go get github.com/jacovdloo/rsaid
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/jacovdloo/rsaid"
)

func main() {

	idNum, err := rsaid.Parse("9506245120081")

	if err != nil {
		fmt.Printf("Invalid South African ID number: %s\n", err)
	} else {
		fmt.Println("Value:", idNum.Value)                                 // Value: 9506245120081
		fmt.Println("DOB:", idNum.DateOfBirth)                             // DOB: 1995-06-24 00:00:00 +0200 SAST
		fmt.Println("Male:", idNum.Gender == rsaid.GenderMale)             // Male: true
		fmt.Println("Citizen:", idNum.Citizenship == rsaid.CitizenshipCitizen) // Citizen: true
	}
}
```
