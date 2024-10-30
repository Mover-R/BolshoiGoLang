package fileutils

import (
	"BolshiGoLang/internal/pkg/storage"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func DataStorageFilePath() string {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("error directory")
		return ""
	}
	return path.Join(cwd, "data.json")
}

func DataStorageFileRead() (*storage.Storage, error) {
	filepath := DataStorageFilePath()
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		emptystorage, err := storage.NewStorage()
		if err != nil {
			return nil, fmt.Errorf("storage error: %w", err)
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

func DataStorageFileWrite(s *storage.Storage) error {
	tempFile, err := os.CreateTemp(filepath.Dir("data.json"), "temp-")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("fail marshal storage: %w", err)
	}
	if _, err := tempFile.Write(data); err != nil {
		return fmt.Errorf("fail writing data: %w", err)
	}
	tempFile.Close()

	return os.Rename(tempFile.Name(), "data.json")
}
