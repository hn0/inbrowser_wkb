package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"model"
	"os"
	"path/filepath"
	"strings"
)

type DB struct {
	file  string
	table string
	conn  *sql.DB
}

func GetConn(file string) *DB {

	// check file signature
	//  refer: https://www.sqlite.org/fileformat.html
	b := make([]byte, 16)
	if fp, err := os.Open(file); err == nil {
		fp.Read(b)
		if string(b[:]) == "SQLite format 3\000" {
			conn, _ := sql.Open("sqlite3", file)
			_, name := filepath.Split(file)
			db := DB{
				file,
				strings.Replace(name, ".sqlite", "", 1),
				conn,
			}
			return &db
		}
	}

	return nil
}

func (db *DB) GetMetadata() {
	// defer db.conn.Close()
	names := []string{"ogc_fid", "statefp", "countryfp"}
	fields := model.CreateFields(names)

	fmt.Println(fields)
	fmt.Println(fields.GetColumns(", "))
	// db.execQuery(fields)
}

// func (db *DB) execQuery(fields []string) {
// 	q := fmt.Sprintf("SELECT %s FROM %s", strings.Join(fields, ","), db.table)
// 	fmt.Println(q)
// 	fmt.Println("continue with execution of the query!@")
// }

func (db *DB) GetSource() string {
	return db.file
}
