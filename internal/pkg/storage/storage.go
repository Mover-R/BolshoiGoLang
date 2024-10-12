package storage

import (
	"strconv"

	"go.uber.org/zap"
)

type Value struct {
	s    string
	kind string
}

type Storage struct {
	inner  map[string]Value
	logger *zap.Logger
}

func NewStorage() (Storage, error) {
	logger, err := zap.NewProduction()

	if err != nil {
		return Storage{}, err
	}

	defer logger.Sync()
	logger.Info("created new storage")

	return Storage{
		inner:  make(map[string]Value),
		logger: logger,
	}, nil
}

func (r Storage) Set(key, value string) {
	switch GetType(value) {
	case "D":
		r.inner[key] = Value{s: value, kind: "D"}
	case "Fl64":
		r.inner[key] = Value{s: value, kind: "Fl64"}
	case "S":
		r.inner[key] = Value{s: value, kind: "S"}
	}

	r.logger.Info("key set")
	r.logger.Sync()
}

func (r Storage) Get(key string) *string {
	res, ok := r.inner[key]
	if !ok {
		return nil
	}

	return &res.s
}

func GetType(value string) string {
	if _, err := strconv.Atoi(value); err == nil {
		return "D"
	}
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return "Fl64"
	}
	return "S"
}
