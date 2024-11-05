package main

import (
	"BolshiGoLang/fileutils"
	"BolshiGoLang/internal/pkg/server"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	r, err := fileutils.DataStorageFileRead()
	if err != nil {
		panic(err)
	}

	port:="8090" //os.Getenv("BASIC_SERVER_PORT")
	s := server.NewServer(":" + port, r)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s.Start()
	}()

	<-signalChan

	err = fileutils.DataStorageFileWrite(r)
	if err != nil {
		fmt.Println("FAILED")
		return
	}
}
