// Data base
package storage

import (
	"go.uber.org/zap"
)

type Storage struct {
	inner  map[string]string
	logger *zap.Logger
}

func InitStorage() (Storage, error) {
	logger, err := zap.NewProduction()

	if err != nil {
		return Storage{}, err
	}
	defer logger.Sync()
	logger.Info("Created new storage!")

	return Storage{
		inner:  make(map[string]string),
		logger: logger,
	}, nil
}

func (r Storage) Set(key, value string) {
	r.inner[key] = value

	r.logger.Info("Key set")
	r.logger.Sync()
}

func (r Storage) Get(key string) *string {
	res, ok := r.inner[key]
	if !ok {
		return nil
	}
	r.logger.Info("key obtained")
	r.logger.Sync()
	return &res
}
