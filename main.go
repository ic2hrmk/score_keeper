package main

import (
	"flag"
	"score_keeper/server"
	"score_keeper/db"
	"log"
)

var port string
var dbFileName string

func init() {
	flag.StringVar(&port, "port", "8080", "port to run on")
	flag.StringVar(&dbFileName, "db", "score.sqlite3", "DB file to communicate with")

	flag.Parse()
}

func main() {
	dbInstance, err := db.NewDataBase(dbFileName)
	if err != nil {
		log.Fatal(err.Error())
	}

	serverInstance := server.NewServer(port, dbInstance)
	serverInstance.Run()
}