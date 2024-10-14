package storage

import (
	"fmt"
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

const (
	KindDigit   = "D"
	KindFloat64 = "Fl64"
	KindString  = "S"
)

func NewStorage() (Storage, error) {
	logger, err := zap.NewProduction(zap.IncreaseLevel(zap.FatalLevel))

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
	case KindDigit:
		r.inner[key] = Value{s: value, kind: "D"}
	case KindFloat64:
		r.inner[key] = Value{s: value, kind: "Fl64"}
	case KindString:
		r.inner[key] = Value{s: value, kind: "S"}
	}

	r.logger.Info("key set")
}

func (r Storage) Get(key string) (string, error) {
	res, ok := r.inner[key]
	if !ok {
		return "", fmt.Errorf("No such key: %q", key)
	}
	return res.s, nil
}

func GetType(value string) string {
	if _, err := strconv.Atoi(value); err == nil {
		return KindDigit
	}
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return KindFloat64
	}
	return KindString
}
