package fileutils

import (
	"BolshiGoLang/internal/pkg/storage"
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func FilePath() string {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("error directory")
		return ""
	}
	return path.Join(cwd, "data.json")
}

func FileRead() (*storage.Storage, error) {
	filepath := FilePath()
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		fmt.Println("Creating file...")

		file, err := os.Create(filepath)
		if err != nil {
			fmt.Println("Fail creating file")
			return nil, err
		}
		defer file.Close()

		emptystorage, err := storage.NewStorage()
		if err != nil {
			fmt.Println("storage error")
			return nil, err
		}
		data, err := json.Marshal(emptystorage)
		if err != nil {
			fmt.Println("fail marshal storage")
			return nil, err
		}

		_, err = file.Write(data)
		if err != nil {
			fmt.Println("Fail writing file")
			return nil, err
		}

		return &emptystorage, nil
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("fail reading file")
		return nil, err
	}
	var s storage.Storage
	err = json.Unmarshal(data, &s)
	if err != nil {
		fmt.Println("fail unmarshal file data")
		return nil, err
	}

	return &s, nil
}

func FileWrite(s *storage.Storage) error {
	filepath := FilePath()
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println("fail creating file")
		return err
	}
	defer file.Close()

	data, err := json.Marshal(s)
	if err != nil {
		fmt.Println("fail marshal storage")
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		fmt.Println("fail writing data")
		return err
	}
	return nil
}
