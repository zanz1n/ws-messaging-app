package services

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/zanz1n/ws-messaging-app/internal/dba"
)

func NewDbProvider() (*dba.Queries, *sql.DB) {
	uri := ConfigProvider().DatabaseUri

	conn, err := sql.Open("postgres", uri)

	if err != nil {
		panic(err)
	}

	db := dba.New(conn)

	return db, conn
}
