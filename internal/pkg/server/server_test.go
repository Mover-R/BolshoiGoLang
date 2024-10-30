package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Request struct {
	Value string
}

type Response struct {
	Value string `json:"value"`
}

type testCase struct {
	name  string
	key   string
	value string
	kind  string
}

var tests = []testCase{
	{"Test1", "key1", "123", "D"},
	{"Test2", "key2", "123.45", "Fl64"},
	{"Test3", "key3", "hello", "S"},
	{"Test4", "key4", "0", "D"},
	{"Test5", "key5", "-42", "D"},
	{"Test6", "key6", "", "S"},
}

func TestHandlerHealth(t *testing.T) {
	t.Run("Health", func(t *testing.T) {
		cli := http.Client{}

		resp, err := cli.Get("http://localhost:8090/health")
		if err != nil {
			t.Error("no response health")
			return
		}

		assert.Equal(t, resp.StatusCode, http.StatusOK)
	})
}

func TestHandlerSet(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var body Request
			body = Request{
				Value: test.value,
			}
			bodyJSON, err := json.Marshal(&body)
			if err != nil {
				t.Error()
				return
			}

			cli := http.Client{}

			host := "http://localhost:8090/scalar/set/" + test.key

			req, err := http.NewRequest(http.MethodPut, host, bytes.NewBuffer(bodyJSON))
			if err != nil {
				t.Error("error creating request")
				return
			}

			req.Header.Set("Content-Type", "application/json")

			resp, err := cli.Do(req)
			if err != nil {
				t.Error("error send json")
				return
			}
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}

func TestHandlerGet(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cli := http.Client{}
			host := "http://localhost:8090/scalar/get/" + test.key

			resp, err := cli.Get(host)
			if err != nil {
				t.Error("failed get response")
			}
			defer resp.Body.Close()

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return
			}

			var val Response
			json.Unmarshal(bodyBytes, &val)
			assert.Equal(t, test.value, val.Value)
		})
	}
}
