package storage

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error

	DB, err = sql.Open("sqlite", "lan-share.db")
	if err != nil {
		log.Fatal("failed to open database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	log.Println("SQLite connected")

	initTables()
}

func initTables() {
	createDeviceTable := `
	CREATE TABLE IF NOT EXISTS device_identity (
		id TEXT PRIMARY KEY,
		name TEXT,
		device_type TEXT,
		os TEXT,
		os_version TEXT,
		arch TEXT,
		hostname TEXT,
		public_key TEXT,
		private_key TEXT,
		created_at INTEGER
	);
	`

	_, err := DB.Exec(createDeviceTable)
	if err != nil {
		log.Fatal("failed to create tables:", err)
	}

	log.Println("Tables initialized")
}
