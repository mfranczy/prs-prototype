package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresDB() *sqlx.DB {
	db, err := sqlx.Open("postgres", "user=prsu dbname=prsdb sslmode=disable password='prsu'")
	if err != nil {
		log.Fatal("Error", err)
	}
	return db
}


// add execute method
// to cover postgres execution
