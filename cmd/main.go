package main

import (
	"BolshiGoLang/internal/pkg/storage"
	"fmt"
)

func main() {
	s, err := storage.NewStorage()
	if err != nil {
		fmt.Println("Something broke(")
		return
	}
	s.Set("0", "mamba")
	s.Set("1", "22")
	s.Set("2", "22.2")
	fmt.Println(*s.Get("0"), storage.Get_type(*s.Get("0")))
	fmt.Println(*s.Get("1"), storage.Get_type(*s.Get("1")))
	fmt.Println(*s.Get("2"), storage.Get_type(*s.Get("2")))
}
