package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
)

func initDB() *sql.DB {
	// 1. 사용자 홈 디렉토리 하위에 db 파일 저장 경로 설정
	home, _ := os.UserHomeDir()
	dbDir := filepath.Join(home, ".ok2pus")
	dbPath := filepath.Join(dbDir, "hosts.db")

	// 2. 폴더 없으면 생성
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		os.MkdirAll(dbDir, 0755)
	}

	// 3. sqlite 연결
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// 4. table 초기화
	query := `
	CREATE TABLE IF NOT EXISTS ssh_hosts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		alias TEXT UNIQUE,
		host TEXT,
		user TEXT,
		port INTEGER DEFAULT 22
	);`

	if _, err := db.Exec(query); err != nil {
		log.Fatal("Table Creation Failed. : ", err)
	}

	return db
}