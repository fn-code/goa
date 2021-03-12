package main

import (
	"os"

	"github.com/fn-code/goa"
	"github.com/google/logger"
)

const (
	// address = "0.0.0.0:5005"
	address = "192.168.1.14:5005"
)

var (
	// Lg is use for define logger
	Lg *logger.Logger
)

func main() {
	fl, err := os.OpenFile("./logger/log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("failed to open file: %v\n", err)
	}

	// Set logger verbose to false if don't want logger to show in stdout
	Lg = logger.Init("AudioClient", true, true, fl)
	defer Lg.Close()

	conn, err := goa.Dial(address)
	if err != nil {
		logger.Fatal(err)
	}
	defer conn.Close()
	conn.SendAudio()

}
