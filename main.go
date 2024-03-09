package main

import (
	"os"
	"time"

	"github.com/pyxvlad/proiect-ipdp/routes"
	"github.com/rs/zerolog"
)

func main() {

	logFile, err := os.Create("ipdp.log")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := logFile.Close()
		if err != nil {
			panic(err)
		}
	}()
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}
	writer := zerolog.MultiLevelWriter(logFile, consoleWriter)
	log := zerolog.New(writer).With().Timestamp().Logger()

	routes.ListenAndServe(&log)
}
