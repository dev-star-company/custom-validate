package validate

import (
	"testing"
)

type TestStruct struct {
	Value int
}

func TestValidate_IssueZeroValue(t *testing.T) {
	fields := map[string]string{
		"Value": "required,numeric,gte=0",
	}

	// Test with 0
	in := TestStruct{Value: 0}
	err := Validate(fields, in)
	if err != nil {
		t.Errorf("Expected nil error for Value=0, got: %v", err)
	}

	// Test with pointer to struct with 0
	err = Validate(fields, &in)
	if err != nil {
		t.Errorf("Expected nil error for Value=0 (pointer), got: %v", err)
	}
}

func TestValidate_Map(t *testing.T) {
	fields := map[string]string{
		"Value": "required,numeric,gte=0",
	}

	// Test with map containing 0
	inMap := map[string]int{"Value": 0}
	err := Validate(fields, inMap)
	if err != nil {
		t.Errorf("Expected nil error for Map Value=0, got: %v", err)
	}

	// Test with map containing pointer to 0
	zero := 0
	inMapPointer := map[string]*int{"Value": &zero}
	err = Validate(fields, inMapPointer)
	if err != nil {
		t.Errorf("Expected nil error for Map Value=&0, got: %v", err)
	}
}

func TestValidate_DoesNotExist(t *testing.T) {
	// Case 1: Struct without the field
	type EmptyStruct struct{}
	err := Validate(map[string]string{"Value": "required"}, EmptyStruct{})
	if err == nil {
		t.Error("Expected error for missing field in struct, got nil")
	} else {
		t.Logf("Expected error caught: %v", err)
	}

	// Case 2: Map without the field
	err = Validate(map[string]string{"Value": "required"}, map[string]int{"Other": 1})
	if err == nil {
		t.Error("Expected error for missing field in map, got nil")
	} else {
		t.Logf("Expected error caught: %v", err)
	}
}

func TestValidate_RequiredStillWorks(t *testing.T) {
	fields := map[string]string{
		"Name": "required",
	}

	// For strings, "" is zero value.
	// Our workaround only targets numeric 0.
	err := Validate(fields, map[string]string{"Name": ""})
	if err == nil {
		t.Error("Expected error for empty string Name with required, got nil")
	}
}
