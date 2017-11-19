package db

import (
	"database/sql"
	_ "fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type DB struct {
	file string
	conn *sql.DB
}

func GetConn(file string) *DB {

	// check file signature
	//  refer: https://www.sqlite.org/fileformat.html
	b := make([]byte, 16)
	if fp, err := os.Open(file); err == nil {
		fp.Read(b)
		if string(b[:]) == "SQLite format 3\000" {
			conn, _ := sql.Open("sqlite3", file)
			db := DB{
				file,
				conn,
			}
			return &db
		}
	}

	return nil
}

func (db *DB) GetSource() string {
	return db.file
}
