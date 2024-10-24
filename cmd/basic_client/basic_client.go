package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type x2Request struct {
	Val int
}

type x2Response struct {
	Val int
}

func cusnsumeResponse(r io.ReadCloser) {
	_, err := io.Copy(io.Discard, r)
	if err != nil {
		log.Print(err)
	}

	if err := r.Close(); err != nil {
		log.Print(err)
	}
}

func main() {
	cli := http.Client{}

	resp, err := cli.Get("http://127.0.0.1:8090/health")
	if err != nil {
		log.Fatal()
	}

	defer cusnsumeResponse(resp.Body)

	fmt.Println(resp.StatusCode)

	req := x2Request{
		Val: 2,
	}
	body, err := json.Marshal(&req)
	if err != nil {
		log.Fatal(err)
	}
	resp, err = cli.Post("http://127.0.0.1:8090/x2", "", bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	var x2Resp x2Response
	if err = json.NewDecoder(resp.Body).Decode(&x2Resp); err != nil {
		log.Fatal(err)
	}

	fmt.Println(x2Resp.Val)

	defer cusnsumeResponse(resp.Body)

	fmt.Println(resp.StatusCode)
}
