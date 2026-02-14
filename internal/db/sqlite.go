package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitDB() *sql.DB {
	home, _ := os.UserHomeDir()
	dbDir := filepath.Join(home, ".ok2pus")
	dbPath := filepath.Join(dbDir, "hosts.db")

	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		os.MkdirAll(dbDir, 0o755)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS ssh_hosts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		alias TEXT UNIQUE,
		host TEXT,
		user TEXT,
		port INTEGER DEFAULT 22,

		auth_type TEXT DEFAULT "Password",
		key_path TEXT
	);`

	if _, err := db.Exec(query); err != nil {
		log.Fatal("Table Creation Failed. : ", err)
	}

	return db
}

func DropDB(db *sql.DB) {
	home, _ := os.UserHomeDir()
	dbDir := filepath.Join(home, ".ok2pus")
	dbPath := filepath.Join(dbDir, "hosts.db")

	err := os.Remove(dbPath)

	if err != nil {
		fmt.Printf("Failed to delete DB file: %v\n", err)
	} else {
		fmt.Printf("Successfully deleted Database file.")
	}

	*db = *InitDB()
	fmt.Println("New database file has been initialized.")
}
