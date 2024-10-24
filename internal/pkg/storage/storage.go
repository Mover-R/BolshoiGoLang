package storage

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Value struct {
	S    string `json:"s"`
	Kind string `json:"kind"`
}

type Storage struct {
	Inner        map[string]Value             `json:"inner"`
	ArrayStore   map[string][]string          `json:"arrayStore"`
	Dictionary   map[string]map[string]string `json:"dict"`
	ExperationAt map[string]time.Time         `json:"experationAt"`
	logger       *zap.Logger
	Mu           sync.RWMutex
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
	storage := Storage{
		Inner:        make(map[string]Value),
		ArrayStore:   make(map[string][]string),
		Dictionary:   make(map[string]map[string]string),
		ExperationAt: make(map[string]time.Time),
		logger:       logger,
	}

	go storage.startCleaning()

	return storage, nil
}

func (r *Storage) startCleaning() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		r.cleanExpiredKeys()
	}
}

func (r *Storage) cleanExpiredKeys() {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	for key, exp := range r.ExperationAt {
		if time.Now().After(exp) {
			delete(r.ExperationAt, key)
			delete(r.Inner, key)
		}
	}
}

func (r *Storage) Set(key, value string, exp ...int) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	var ex int = 0
	if len(exp) > 0 {
		ex = exp[0]
	}

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

	if ex > 0 {
		r.ExperationAt[key] = time.Now().Add(time.Duration(ex) * time.Second)
	} else {
		delete(r.ExperationAt, key)
	}

	r.logger.Info("key set")
}

func (r *Storage) Get(key string) (string, error) {
	r.Mu.RLock()
	defer r.Mu.RUnlock()

	if r.Inner == nil {
		r.Inner = make(map[string]Value)
	}
	if r.logger == nil {
		r.logger, _ = zap.NewProduction(zap.IncreaseLevel(zap.FatalLevel))
	}

	if exp, ok := r.ExperationAt[key]; ok && time.Now().After(exp) {
		r.Mu.Unlock()
		r.Mu.Lock()
		delete(r.ExperationAt, key)
		delete(r.Inner, key)
		r.Mu.Unlock()
		return "", errors.New("no such key")
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

func (r *Storage) PrintArr(key string) {
	if _, ok := r.ArrayStore[key]; !ok {
		return
	}
	fmt.Println(r.ArrayStore[key])
}

func (r *Storage) LPUSH(key string, elements ...string) {
	r.ArrayStore[key] = append(elements, r.ArrayStore[key]...)
}

func (r *Storage) LPOP(key string, count ...int) ([]string, error) {
	if _, exists := r.ArrayStore[key]; !exists {
		return []string{}, errors.New("no such key")
	}

	list := r.ArrayStore[key]
	to_return := []string{}
	var left, right int
	if len(count) == 0 {
		left, right = 0, 1
	} else if len(count) == 1 {
		left, right = 0, count[0]
	} else if len(count) == 2 {
		left, right = count[0], count[1]
		if left < 0 {
			left += len(list)
		}
		if right < 0 {
			right += len(list)
		}
	} else {
		return []string{}, errors.New("incorrect data")
	}

	if left < 0 || left > len(list) || right < left || right > len(list) {
		return []string{}, errors.New("incorrect data")
	}

	res := list[left : right+1]
	to_return = append(to_return, res...)
	list = append(list[:left], list[right+1:]...)
	r.ArrayStore[key] = list

	return to_return, nil
}

func (r *Storage) RPUSH(key string, elements ...string) {
	r.ArrayStore[key] = append(r.ArrayStore[key], elements...)
}

func (r *Storage) RPOP(key string, count ...int) ([]string, error) {
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
		slices.Reverse(res)
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
		slices.Reverse(res)
		to_return = append(to_return, res...)
		list = append(list[:count[1]], list[1+count[0]:]...)
		r.ArrayStore[key] = list
		return to_return, nil
	}

	return []string{}, errors.New("count = {'', '1', '2'}")
}

func (r *Storage) RADDTOSET(key string, elements ...string) {
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

func (r *Storage) LSET(key string, index int, element string) error {
	list, ok := r.ArrayStore[key]
	if !ok {
		return fmt.Errorf("no such key")
	}

	if len(list) < index || index < 0 {
		return fmt.Errorf("out of range")
	}

	list[index-1] = element
	r.logger.Info("Value changed")
	return nil
}

func (r *Storage) LGET(key string, index int) (string, error) {
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

type entry struct {
	field string
	val   string
}

func (r *Storage) HSET(key string, data ...entry) {
	if r.Dictionary[key] == nil {
		r.Dictionary[key] = make(map[string]string)
	}

	for _, el := range data {
		r.Dictionary[key][el.field] = el.val
	}
}
