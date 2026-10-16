package main

import (
	"database/sql"
	_ "embed"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schema string

func openDB() *sql.DB {
	db, err := sql.Open("sqlite", "file:music.db?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(schema); err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	db := openDB()
	defer db.Close()

	if err := scanLibrary(db, "/home/alice/Music"); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	registerRoutes(mux, db)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
