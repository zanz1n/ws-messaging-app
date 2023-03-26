package services

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	"github.com/zanz1n/ws-messaging-app/internal/dba"
)

func NewDbProvider() (*dba.Queries, *sql.DB) {
	uri := os.Getenv("DATABASE_URI")

	conn, err := sql.Open("postgres", uri)

	if err != nil {
		panic(err)
	}

	db := dba.New(conn)

	return db, conn
}
