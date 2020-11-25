package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Database struct {
	Conn *sql.DB
}

func Initialize(host string, port int, username, password, database string) (Database, error) {
	db := Database{}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, database)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return db, fmt.Errorf("opening database: %v", err)
	}

	db.Conn = conn
	err = db.Conn.Ping()
	if err != nil {
		return db, fmt.Errorf("pinging database: %v", err)
	}
	return db, nil
}
