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

	dbPath, found := os.LookupEnv("IPDP_DB")
	if !found {
		dbPath = "./ipdp.db"
	}
	log.Info().Msgf("Opening DB at %s", dbPath)

	sqliteDB, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		log.Fatal().Err(err).Msg("While trying to open database")
	}

	err = database.MigrateDB(sqliteDB)
	if err != nil {
		log.Fatal().Err(err).Msg("While trying to migrate database")
	}

	coverPath, found := os.LookupEnv("IPDP_COVERPATH")
	if !found {
		coverPath = "./covers"
	}
	log.Info().Msgf("Storing covers at %s", coverPath)

	err = os.MkdirAll(coverPath, 0777)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't create folder for storing covers")
	}

	err = os.MkdirAll("/tmp/ipdp-img/", 0777)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't create temp folder for storing covers")
	}

	routes.ListenAndServe(&log, sqliteDB, coverPath)
}
