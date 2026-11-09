package user

import (
	"encoding/json"
	"testing"
)

// Decoding into a reused User must not keep the keys of an earlier decode in Raw.
func TestUnmarshalJSONReplacesRaw(t *testing.T) {
	var u User
	if err := json.Unmarshal([]byte(`{"id":1,"favourite_colour":"blue"}`), &u); err != nil {
		t.Fatalf("first Unmarshal() error = %v", err)
	}
	if err := json.Unmarshal([]byte(`{"id":2}`), &u); err != nil {
		t.Fatalf("second Unmarshal() error = %v", err)
	}
	var colour string
	if ok, err := u.FieldValue("favourite_colour", &colour); ok || err != nil {
		t.Errorf("FieldValue(favourite_colour) = %v, %v, %q; want false, nil: the key came from the earlier decode", ok, err, colour)
	}
}
