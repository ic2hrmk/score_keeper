package server

import (
	"score_keeper/db"

	"errors"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

const (
	newRecordPath = "/new_record"
	leaderBoardPath = "/leader_board"
)

const (
	logo = `
		   _____                      _  __
		  / ____|                    | |/ /
		 | (___   ___ ___  _ __ ___  | ' / ___  ___ _ __   ___ _ __
		  \___ \ / __/ _ \| '__/ _ \ |  < / _ \/ _ \ '_ \ / _ \ '__|
		  ____) | (_| (_) | | |  __/ | . \  __/  __/ |_) |  __/ |
		 |_____/ \___\___/|_|  \___| |_|\_\___|\___| .__/ \___|_|
							   | |
							   |_|
	`
)

type Server struct {
	port string
	db *db.DataBase
}

func NewServer(port string, db *db.DataBase) *Server {
	return &Server{
		port: port,
		db: db,
	}
}

func (server *Server) getLeaderBoard(rw http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		leaderBoardJSON, err := server.db.GetLeaderBoardJSON()
		if err != nil {
			log.Println("ERROR: " + err.Error())
		}

		rw.Write(leaderBoardJSON)

		log.Println(leaderBoardPath + ": GET leader board")
	} else {
		http.Error(rw,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)

		log.Println(fmt.Sprintf("method %s instead of %s", req.Method, http.MethodGet))
	}
}

func (server *Server) saveNewRecord(rw http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println("ERROR: " + err.Error())
		}

		var record *db.Record
		err = json.Unmarshal(body, &record)
		if err != nil {
			log.Println("ERROR: " + err.Error())
		}

		if record == nil {
			http.Error(rw,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)

			log.Println("unrecognised JSON object")

			return
		}

		err = server.db.SaveRecord(record)
		if err != nil {
			log.Println("ERROR: " + err.Error())
		}

		recordJSON, _ := json.MarshalIndent(record, "  ", "  ")
		rw.Write(recordJSON)

		log.Println(leaderBoardPath + ": POST new record")
	} else {
		http.Error(rw,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)

		log.Println(fmt.Sprintf("method %s instead of %s", req.Method, http.MethodPost))
	}
}

func (server *Server) Run() {
	var err error

	if server.port == "" {
		err = errors.New("port not set")
		return
	}

	http.HandleFunc(newRecordPath, server.saveNewRecord)
	http.HandleFunc(leaderBoardPath, server.getLeaderBoard)

	log.Println(logo)
	log.Println("Score keeper server started at port " + server.port +". Enjoy your game!")
	log.Println("Awaiting for new records...")

	err = http.ListenAndServe(":" + server.port, nil)
	log.Fatal(err.Error())

}