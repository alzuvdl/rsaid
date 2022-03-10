# 🇿🇦 South African ID Validator for Go

> Inspired by https://www.npmjs.com/package/south-african-id-validator

Validate South African ID numbers. Takes eligibility age into account (16 years).

The following information can be derived if the ID Number is valid:

- Gender
- Citizenship
- Date of birth

## Details

A South African ID number is a 13-digit number with the following format: `YYMMDDSNNNCAZ`.

The first six digits `YYMMDD` indicates a person's birth date. For example, 24 June 1995 becomes **950624**. Yes, RWC winners! 🏉

Although rare, it can happen that someone’s birth date does not correspond with their ID number.

> {950624} SNNNCAZ.


The next digit `S` indicates a person's gender. Females have a number of `0 to 4`, while males are `5 to 9`.

> {950624} {5} NNNCAZ

The next three digits `NNN` indicates the position of birth by registry based on the date of birth and gender.

If for example, the number is 120, it means that person was the 120th person to be registered as having been born on that particular day, for that gender.

> {950624} {5} {120} CAZ

The next digit `C` indicates citizenship. `0` if the person is a South African citizen, or `1` if the person is a permanent resident.

> {950624} {5} {120} {0} AZ

The next digit `A` was used until the late 1980s to indicate a person’s race. This has been eliminated and
old ID numbers were reissued to remove this.

> {950624} {5} {120} {0} {0} Z

The last digit `Z` is a checksum digit, used to check that the number sequence is accurate using the [Luhn](https://en.wikipedia.org/wiki/Luhn_algorithm) algorithm.

So, ID number `9506245120008` will reflect as the `120th` `male` South African `citizen` born/registered on the `24th of June, 1995`

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

  id := "9506245120008"
  person, err := rsaid.Parse(id)
  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println("Gender:", person.Gender) // Gender: male
    fmt.Println("Citizen:", person.Citizen) // Citizen: true
    fmt.Println("DOB:", person.DOB) //DOB: 1995-06-24 00:00:00 +0200 SAST
  }

  gender, g_err := rsaid.Gender(id)
  if g_err != nil {
    fmt.Println(g_err)
  }

  citizen, c_err := rsaid.IsCitizen(id)
  if c_err != nil {
    fmt.Println(c_err)
  }

  dob, d_err := rsaid.BirthDate(id)
  if d_err != nil {
    fmt.Println(d_err)
  }

  fmt.Printf("Gender: %s, Citizen: %t, DOB: %s\n", gender, citizen, dob)
  // Gender: male, Citizen: true, DOB: 1995-06-24 00:00:00 +0200 SAST
}
```
