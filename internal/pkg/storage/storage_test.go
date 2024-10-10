package storage

import (
	"fmt"
	"testing"
)

func TestGetType(t *testing.T) {
	fmt.Println("Tests for GetType()")
	tests := []struct {
		input    string
		expected string
	}{
		{"123", "D"},
		{"123.45", "Fl64"},
		{"hello", "S"},
		{"0", "D"},
		{"-42", "D"},
	}
	for _, test := range tests {
		res := GetType(test.input)
		if res != test.expected {
			t.Errorf("GetType(%q) = %q != Expexted_type %q\n", test.input, res, test.expected)
		}
	}
}
