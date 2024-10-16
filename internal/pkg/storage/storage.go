package storage

import (
	"errors"
	"fmt"
	"strconv"

	"go.uber.org/zap"
)

type Value struct {
	S    string `json:"s"`
	Kind string `json:"kind"`
}

type Storage struct {
	Inner      map[string]Value `json:"inner"`
	logger     *zap.Logger
	ArrayStore map[string][]string `json:"ArrayStore"`
}

const (
	KindDigit   = "D"
	KindFloat64 = "Fl64"
	KindString  = "S"
	KindArray   = "Arr"
)

func NewStorage() (Storage, error) {
	logger, err := zap.NewProduction(zap.IncreaseLevel(zap.FatalLevel))

	if err != nil {
		return Storage{}, err
	}

	defer logger.Sync()
	logger.Info("created new storage")

	return Storage{
		Inner:      make(map[string]Value),
		logger:     logger,
		ArrayStore: make(map[string][]string),
	}, nil
}

func (r Storage) Set(key, value string) {
	if r.Inner == nil {
		r.Inner = make(map[string]Value)
	}
	if r.logger == nil {
		r.logger, _ = zap.NewProduction(zap.IncreaseLevel(zap.FatalLevel))
	}
	switch GetType(value) {
	case KindDigit:
		r.Inner[key] = Value{S: value, Kind: "D"}
	case KindFloat64:
		r.Inner[key] = Value{S: value, Kind: "Fl64"}
	case KindString:
		r.Inner[key] = Value{S: value, Kind: "S"}
	}

	r.logger.Info("key set")
}

func (r Storage) Get(key string) (string, error) {
	if r.Inner == nil {
		r.Inner = make(map[string]Value)
	}
	if r.logger == nil {
		r.logger, _ = zap.NewProduction(zap.IncreaseLevel(zap.FatalLevel))
	}
	res, ok := r.Inner[key]
	if !ok {
		return "", errors.New("no such key: ")
	}
	return res.S, nil
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

func reverse(slice []string) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func (r Storage) PrintArr(key string) {
	if _, err := r.ArrayStore[key]; !err {
		fmt.Errorf("no such key")
		return
	}
	fmt.Println(r.ArrayStore[key])
}

func (r Storage) LPUSH(key string, elements ...string) {
	if _, err := r.ArrayStore[key]; !err {
		r.ArrayStore[key] = []string{}
	}
	//reverse(elements)
	r.ArrayStore[key] = append(elements, r.ArrayStore[key]...)
}

func (r Storage) LPOP(key string, count ...int) ([]string, error) {
	if _, exists := r.ArrayStore[key]; !exists {
		return []string{}, errors.New("no such key")
	}
	list := r.ArrayStore[key]
	to_return := []string{}
	if len(count) == 0 {
		if len(list) <= 0 {
			return []string{}, errors.New("no more elements in list")
		}
		res := list[:1]
		to_return = append(to_return, res...)
		list = list[1:]
		r.ArrayStore[key] = list
		return to_return, nil
	} else if len(count) == 1 {
		if len(list) < count[0] {
			return []string{}, errors.New("no more elements in list")
		}
		res := list[:count[0]]
		to_return = append(to_return, res...)
		list = list[count[0]:]
		r.ArrayStore[key] = list
		return to_return, nil
	} else if len(count) == 2 {
		if count[0] < 0 {
			count[0] += len(list)
		}
		if count[1] < 0 {
			count[1] += len(list)
		}
		if len(list) < count[1]-count[0] {
			return []string{}, errors.New("no more elements in list")
		}
		res := list[count[0] : count[1]+1]
		to_return = append(to_return, res...)
		list = append(list[:count[0]], list[1+count[1]:]...)
		r.ArrayStore[key] = list
		return to_return, nil
	}
	return []string{}, errors.New("count = {'', '1', '2'}")
}

func (r Storage) RPUSH(key string, elements ...string) {
	if _, exists := r.ArrayStore[key]; !exists {
		r.ArrayStore[key] = []string{}
	}
	r.ArrayStore[key] = append(r.ArrayStore[key], elements...)
}

func (r Storage) RPOP(key string, count ...int) ([]string, error) {
	if _, exists := r.ArrayStore[key]; !exists {
		return []string{}, errors.New("no such key")
	}
	list := r.ArrayStore[key]
	to_return := []string{}
	len_list := len(list)
	if len(count) == 0 {
		if len_list <= 0 {
			return []string{}, errors.New("no more elements in list")
		}
		res := list[len_list-1:]
		to_return = append(to_return, res...)
		list = list[:len_list-1]
		r.ArrayStore[key] = list
		return to_return, nil
	} else if len(count) == 1 {
		if len_list < count[0] {
			return []string{}, errors.New("no more elements in list")
		}
		res := list[len_list-count[0]:]
		reverse(res)
		to_return = append(to_return, res...)
		list = list[:len_list-count[0]]
		r.ArrayStore[key] = list
		return to_return, nil
	} else if len(count) == 2 {
		if count[0] < 0 {
			count[0] = -count[0]
		} else {
			count[0] = len_list - count[0]
		}
		if count[1] < 0 {
			count[1] = -count[1]
		} else {
			count[1] = len_list - count[0]
		}
		if len(list) < count[0]-count[1] {
			return []string{}, errors.New("no more elements in list")
		}
		res := list[count[1] : count[0]+1]
		reverse(res)
		to_return = append(to_return, res...)
		list = append(list[:count[1]], list[1+count[0]:]...)
		r.ArrayStore[key] = list
		return to_return, nil
	}
	return []string{}, errors.New("count = {'', '1', '2'}")
}

func (r Storage) RADDTOSET(key string, elements ...string) {
	if _, exists := r.ArrayStore[key]; !exists {
		r.ArrayStore[key] = []string{}
	}
	list := r.ArrayStore[key]
	set := make(map[string]bool)
	for _, el := range list {
		set[el] = true
	}
	for _, el := range elements {
		if _, ok := set[el]; !ok {
			list = append(list, el)
			set[el] = true
		}
	}
	r.ArrayStore[key] = list
}

func (r Storage) LSET(key string, index int, element string) {
	if _, exists := r.ArrayStore[key]; !exists {
		fmt.Println("no such key")
		return
	}
	if len(r.ArrayStore[key]) < index || index < 0 {
		fmt.Println("out of range")
		return
	}
	r.ArrayStore[key][index-1] = element
	fmt.Println("OK")
}

func (r Storage) LGET(key string, index int) (string, error) {
	if _, exists := r.ArrayStore[key]; !exists {
		return "", errors.New("no such key")
	}
	if len(r.ArrayStore[key]) < index || index < 0 {
		return "", errors.New("out of range")
	}
	element := r.ArrayStore[key][index-1]
	fmt.Println(element)
	return element, nil
}
