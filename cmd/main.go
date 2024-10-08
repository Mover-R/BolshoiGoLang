package main

import (
	"BolshiGoLang/internal/pkg/storage"
	"fmt"
	"log"
)

func main() {
	s, err := storage.InitStorage()
	if err != nil {
		log.Fatal(err)
	}
	s.Set("key2", "")

	res1 := s.Get("key2")
	res2 := s.Get("key3")
	fmt.Println(res1, res2)
}
