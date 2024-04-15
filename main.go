package main

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pyxvlad/proiect-ipdp/database"
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

	sqliteDB, err := sql.Open("sqlite3", "ipdp.db")

	if err != nil {
		log.Fatal().Err(err).Msg("While trying to open database")
	}

	err = database.MigrateDB(sqliteDB)
	if err != nil {
		log.Fatal().Err(err).Msg("While trying to migrate database")
	}

	routes.ListenAndServe(&log, sqliteDB)
}
