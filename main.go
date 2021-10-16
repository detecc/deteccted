package main

import (
	"github.com/detecc/deteccted/config"
	"github.com/detecc/deteccted/connection"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	config.GetFlags()

	connection.Start()

	<-quitChannel
}
