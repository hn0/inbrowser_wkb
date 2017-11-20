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
	names := []string{"ogc_fid", "statefp"}
	fields := model.CreateFields(names)

	if cnt := db.execSelect(fields); cnt >= 0 {
		fmt.Println(fields)
	}
}

func (db *DB) execSelect(fields *model.Fields) int {
	defer db.conn.Close()

	if err := db.conn.Ping(); err != nil {
		if conn, err2 := sql.Open("sqlite3", db.file); err2 == nil {
			db.conn = conn
		} else {
			return -2
		}
	}

	cnt := 0
	q := fmt.Sprintf("SELECT %s FROM %s", fields.GetColumns(", "), db.table)
	fmt.Println(q)

	// second argument are array of values for constraint condition
	if rows, err := db.conn.Query(q, nil); err == nil {
		defer rows.Close()
		if coltyps, err := rows.ColumnTypes(); err == nil && len(coltyps) > 0 {

			fmt.Println(coltyps)

			for rows.Next() {
				// todo: retrivie interface from the model
				rowvals := make([]interface{}, 3)
				for i := range rowvals {
					var v interface{}
					rowvals[i] = &v
				}
				if err := rows.Scan(rowvals...); err == nil {
					cnt = cnt + 1
				}
			}
			return cnt
		}
	} else {
		fmt.Println(err)
	}

	return -1
}

func (db *DB) GetSource() string {
	return db.file
}
