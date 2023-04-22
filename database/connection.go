package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	connection, err := sql.Open("mysql", "jacob:jacob@tcp(localhost:3306)/CSC4990")
	if err != nil {
		panic("could not connect to the database")
	}
	DB = connection
}
