package db

import (
	_ "github.com/mattn/go-sqlite3"

	"errors"
	"log"
	"database/sql"
	"encoding/json"
)

type DataBase struct {
	fileName   string
	connection *sql.DB
}

func NewDataBase(fileName string) (db *DataBase, err error) {
	db = &DataBase{
		fileName: fileName,
	}

	db.connection, err = sql.Open("sqlite3", db.fileName)
	if err != nil {
		return
	}

	err = db.initDB()

	return
}

func (db *DataBase) initDB() (err error) {
	sqlStmt := `CREATE TABLE IF NOT EXISTS score (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			player     TEXT NULL,
			score      REAL
		    );`

	_, err = db.connection.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	return
}

func (db *DataBase) GetLeaderBoardJSON() (lbJSON []byte, err error) {
	sqlQuery := `SELECT id, player, score
			FROM score
			ORDER BY score DESC`

	rows, err := db.connection.Query(sqlQuery)
	defer rows.Close()

	if err != nil {
		return
	}

	var records []Record

	defer rows.Close()
	for rows.Next() {
		var record Record

		err = rows.Scan(&record.Id, &record.Player, &record.Score)
		if err != nil {
			return
		}

		records = append(records, record)
	}

	err = rows.Err()
	if err != nil {
		return
	}

	lbJSON, err= json.MarshalIndent(records, "  ", "  ")

	return
}

func (db *DataBase) SaveRecord(record *Record) (err error) {
	switch {
	case record == nil:
		err = errors.New("empty record cannot be saved")
	case record.Player == "":
		err = errors.New("player name cannot be empty")
	}

	if err != nil {
		return
	}

	tx, err := db.connection.Begin()
	if err != nil {
		return
	}
	stmt, err := tx.Prepare(`INSERT INTO score(player, score) SELECT $1, $2`)
	if err != nil {
		return
	}

	defer stmt.Close()

	if _, err = stmt.Exec(record.Player, record.Score); err != nil {
		return
	}

	err = tx.Commit()

	return
}
