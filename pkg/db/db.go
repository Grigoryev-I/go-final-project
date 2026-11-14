package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX scheduler_date ON scheduler (date);
`
const defaultDBFile = "scheduler.db"

func GetDBFile() string {
	if envFile := os.Getenv("TODO_DBFILE"); envFile != "" {
		return envFile
	}

	return defaultDBFile
}

var db *sql.DB

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		if _, err = db.Exec(schema); err != nil {
			db.Close()
			return err
		}
	}

	return nil
}

func Close() error {
	return db.Close()
}
