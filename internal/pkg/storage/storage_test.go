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

func TestSetGet(t *testing.T) {
	fmt.Println("Tests for Set() Get()")
	s, err := NewStorage()
	if err != nil {
		t.Fatal("Storage didnt created")
	}
	tests := []struct {
		in  string
		out string
	}{
		{"key1", "value1"},
		{"key2", "value2"},
		{"key3", "value3"},
		{"key4", ""},
	}
	for _, test := range tests {
		s.Set(test.in, test.out)
	}
	for _, test := range tests {
		res := s.Get(test.in)
		if res == nil {
			t.Fatal("No such key")
		}
		if (*res) != test.out {
			t.Errorf("Get(%q) = %q != Expected %q", test.in, *res, test.out)
		}
	}
}
