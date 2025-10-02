package common

import (
	"encoding/json"
	"strconv"
)

// FlexInt is a custom type that can unmarshal both string and number JSON values to int.
// It solves the problem: json: cannot unmarshal number into Go struct field of type string
//
// Example usage:
//
//	type MyStruct struct {
//	    Sort common.FlexInt `json:"sort"`
//	}
//
// This will correctly unmarshal both:
//
//	{"sort": 42}        -> FlexInt(42)
//	{"sort": "42"}      -> FlexInt(42)
//	{"sort": "invalid"} -> FlexInt(0)
type FlexInt int

// UnmarshalJSON implements json.Unmarshaler interface
// Accepts both string and number, returns 0 if conversion fails
func (fi *FlexInt) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as int first
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		*fi = FlexInt(i)
		return nil
	}

	// Try to unmarshal as string
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		if num, err := strconv.Atoi(s); err == nil {
			*fi = FlexInt(num)
			return nil
		}
	}

	// If both fail, set to 0
	*fi = 0
	return nil
}

// MarshalJSON implements json.Marshaler interface
// Always marshals as int
func (fi FlexInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(fi))
}

// Int returns the int value
func (fi FlexInt) Int() int {
	return int(fi)
}

type (
	ContactData struct {
		Id    string  `json:"id"`
		Type  string  `json:"type"`
		Value string  `json:"value"`
		Main  bool    `json:"main"`
		Sort  FlexInt `json:"sort"`
	}

	Meta struct {
		CurrentPage int `json:"current_page"`
		From        int `json:"from"`
		LastPage    int `json:"last_page"`
		PerPage     int `json:"per_page"`
		To          int `json:"to"`
		Total       int `json:"total"`
	}

	Links struct {
		First string `json:"first"`
		Last  string `json:"last"`
		Prev  any    `json:"prev"`
		Next  any    `json:"next"`
	}

	SocialMedia struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Link string `json:"link"`
	}
)
