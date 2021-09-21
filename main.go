package main

import (
	"os"
	"os/signal"
	"github.com/detecc/deteccted/config"
	"github.com/detecc/deteccted/connection"
	"syscall"
)

func main() {
	config.GetFlags()

	connection.Start()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}