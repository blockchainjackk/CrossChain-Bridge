package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
)

func Connect(host, port, user, pass, dbname string) (*sql.DB, error) {
	var psqlInfo string
	if pass == "" {
		psqlInfo = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable",
			host, user, dbname)
	} else {
		psqlInfo = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
			host, user, pass, dbname)
	}

	// Only add port arg for TCP connections since UNIX domain sockets
	// (specified by a "/" prefix) do not have a port.
	if !strings.HasPrefix(host, "/") {
		psqlInfo += fmt.Sprintf(" port=%s", port)
	}

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}
