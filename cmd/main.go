package main

import (
	"BolshiGoLang/fileutils"
	"BolshiGoLang/internal/pkg/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	r, err := fileutils.DataStorageFileRead()
	if err != nil {
		panic(err)
	}
	s := server.NewServer(":8090", r)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s.Start()
	}()

	<-signalChan

	err = fileutils.DataStorageFileWrite(r)
	if err != nil {
		return
	}
}
