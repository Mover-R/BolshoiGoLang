package benchmarkresults

import (
	"BolshiGoLang/internal/pkg/storage"
	"math/rand"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func disableLogging() *zap.Logger {
	logger := zap.NewNop()
	return logger
}

type bench struct {
	name       string
	count_keys int
}

var cases = []bench{
	{"10", 10},
	{"100", 100},
	{"1000", 1000},
	{"10000", 10000},
}

const alf = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const alf_len = len(alf)

func randomstring(len int) string {
	var res strings.Builder
	for i := 0; i < len; i++ {
		indx := rand.Intn(alf_len)
		res.WriteByte(alf[indx])
	}

	return res.String()
}

func BenchmarkGet(b *testing.B) {
	for _, test := range cases {
		b.Run(test.name, func(bb *testing.B) {
			s, err := storage.NewStorage()

			if err != nil {
				b.Error("Bad storage")
			}

			keys := []string{}
			for i := 0; i < test.count_keys; i++ {
				len_key1, len_val1 := rand.Intn(10), rand.Intn(10)
				key1, val1 := randomstring(len_key1), randomstring(len_val1)
				keys = append(keys, key1)
				s.Set(key1, val1)
			}

			bb.ResetTimer()

			for i := 0; i < bb.N; i++ {
				indx := rand.Intn(len(keys))
				key := keys[indx]
				_, ok := s.Get(key)
				if ok == nil {
					bb.Error("Bad Get function")
				}
			}
		})
	}
}

func BenchmarkSet(b *testing.B) {
	for _, test := range cases {
		b.Run(test.name, func(bb *testing.B) {
			s, err := storage.NewStorage()

			if err != nil {
				b.Error("Bad storage")
			}

			keys := []string{}
			values := "val"
			for i := 0; i < test.count_keys; i++ {
				len_key1 := rand.Intn(10)
				key1 := randomstring(len_key1)
				keys = append(keys, key1)
			}

			bb.ResetTimer()

			for i := 0; i < bb.N; i++ {
				indx := rand.Intn(len(keys))
				s.Set(keys[indx], values)
			}
		})
	}
}

func BenchmarkSetGet(b *testing.B) {
	for _, test := range cases {
		b.Run(test.name, func(bb *testing.B) {
			s, err := storage.NewStorage()

			if err != nil {
				b.Error("Bad storage")
			}

			keys := []string{}
			values := []string{}
			for i := 0; i < test.count_keys; i++ {
				len_key1, len_val1 := rand.Intn(10), rand.Intn(10)
				key1, val1 := randomstring(len_key1), randomstring(len_val1)
				keys = append(keys, key1)
				values = append(values, val1)
			}

			bb.ResetTimer()

			for i := 0; i < bb.N; i++ {
				indx := rand.Intn(len(keys))
				key := keys[indx]
				s.Set(key, values[indx])
				_, ok := s.Get(key)
				if ok == nil {
					bb.Error("Bad Get function")
				}
			}
		})
	}
}
