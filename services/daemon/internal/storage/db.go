package storage

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"

	db "lan-share/daemon/internal/storage/db"
)

var Conn *sql.DB
var Queries *db.Queries

func InitDB() {
	var err error

	Conn, err = sql.Open("sqlite", "lan-share.db")
	if err != nil {
		log.Fatal(err)
	}

	if err = Conn.Ping(); err != nil {
		log.Fatal(err)
	}

	Queries = db.New(Conn)

	log.Println("SQLite + sqlc initialized")

	initTables()
}

func initTables() {
	_, err := Conn.Exec(`
	CREATE TABLE IF NOT EXISTS device_identity (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		device_type TEXT NOT NULL,
		os TEXT NOT NULL,
		os_version TEXT NOT NULL,
		arch TEXT NOT NULL,
		hostname TEXT NOT NULL,
		public_key TEXT NOT NULL,
		private_key TEXT NOT NULL,
		created_at INTEGER NOT NULL
	);
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = Conn.Exec(`
	CREATE TABLE IF NOT EXISTS devices (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		public_key TEXT NOT NULL,
		state TEXT NOT NULL,
		last_seen INTEGER NOT NULL,
		trusted_at INTEGER NOT NULL
	);
	`)
	if err != nil {
		log.Fatal(err)
	}
}
