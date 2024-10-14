package main

import (
	"BolshiGoLang/internal/pkg/storage"
	"fmt"
)

func main() {
	s, err := storage.NewStorage()

	if err != nil {
		fmt.Println("Something broken(")
		return
	}

	s.Set("0", "mamba")
	s.Set("1", "22")
	s.Set("2", "22.2")

	res, ok := s.Get("0")
	if ok != nil {
		fmt.Errorf("Bad storage: no such key: %q", "0")
	}
	res1, ok1 := s.Get("1")
	if ok1 != nil {
		fmt.Errorf("Bad storage: no such key: %q", "0")
	}
	res2, ok2 := s.Get("0")
	if ok2 != nil {
		fmt.Errorf("Bad storage: no such key: %q", "0")
	}
	fmt.Println(res, storage.GetType(res))
	fmt.Println(res1, storage.GetType(res1))
	fmt.Println(res2, storage.GetType(res2))
}
