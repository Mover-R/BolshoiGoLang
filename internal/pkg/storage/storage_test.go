package storage

import (
	"fmt"
	"testing"
)

type testCase struct {
	name  string
	key   string
	value string
	kind  string
}

func TestSetGet(t *testing.T) {
	fmt.Println("Tests for Set(), Get()")
	s, err := NewStorage()
	if err != nil {
		t.Error("Bad storage")
	}
	tests := []testCase{
		{"Test1", "key1", "123", "D"},
		{"Test2", "key2", "123.45", "Fl64"},
		{"Test3", "key3", "hello", "S"},
		{"Test4", "key4", "0", "D"},
		{"Test5", "key5", "-42", "D"},
		{"Test6", "key6", "", "S"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.Set(test.key, test.value)

			res, ok := s.Get(test.key)
			if ok != nil {
				t.Errorf("Bad storage: no such key: %q", test.key)
			}

			if res != test.value {
				t.Error("Values arent equal")
			}
		})
	}
}

func TestSetGetWithType(t *testing.T) {
	fmt.Println("Tests for GetType(), Set, Get")
	s, err := NewStorage()
	if err != nil {
		t.Error("Bad storage")
	}
	tests := []testCase{
		{"Test1", "key1", "123", "D"},
		{"Test2", "key2", "123.45", "Fl64"},
		{"Test3", "key3", "hello", "S"},
		{"Test4", "key4", "0", "D"},
		{"Test5", "key5", "-42", "D"},
		{"Test6", "key6", "", "S"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.Set(test.key, test.value)

			res, ok := s.Get(test.key)
			if ok != nil {
				t.Errorf("Bad storage: no such key: %q", test.key)
			}

			if res != test.value {
				t.Error("Values arent equal")
			}

			if GetType(res) != test.kind {
				t.Error("GetType is not equal to test's type")
			}
		})
	}
}
