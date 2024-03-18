package main

import (
	"os"
	"time"

	"github.com/pyxvlad/proiect-ipdp/models"
	"github.com/pyxvlad/proiect-ipdp/routes"
	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

	sqliteDB := sqlite.Open("ipdp.db")

	db, err := gorm.Open(sqliteDB, &gorm.Config{})

	if err != nil {
		log.Fatal().Err(err).Msg("While trying to open database")
	}

	err = models.AutoMigrate(db)
	if err != nil {
		log.Fatal().Err(err).Msg("While trying to migrate database")
	}

	routes.ListenAndServe(&log, db)
}
