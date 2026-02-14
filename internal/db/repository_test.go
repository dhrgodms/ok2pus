package db

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"

	"ok2pus/internal/model"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
	}

	t.Cleanup(func() { db.Close() })
	return db
}

func newTestHost(alias string) model.SSHHost {
	return model.SSHHost{
		Alias:    alias,
		Host:     "192.168.1.1",
		User:     "root",
		Port:     22,
		AuthType: "Password",
	}
}

// --- AddHost ---

func TestAddHost_Success(t *testing.T) {
	db := setupTestDB(t)

	err := AddHost(db, newTestHost("myserver"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestAddHost_DuplicateAlias(t *testing.T) {
	db := setupTestDB(t)

	_ = AddHost(db, newTestHost("dup"))
	err := AddHost(db, newTestHost("dup"))

	if err == nil {
		t.Fatal("expected error for duplicate alias, got nil")
	}
}

// --- GetAllHost ---

func TestGetAllHost_Empty(t *testing.T) {
	db := setupTestDB(t)

	hosts, err := GetAllHost(db)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(hosts) != 0 {
		t.Fatalf("expected 0 hosts, got %d", len(hosts))
	}
}

func TestGetAllHost_WithData(t *testing.T) {
	db := setupTestDB(t)
	_ = AddHost(db, newTestHost("server1"))
	_ = AddHost(db, newTestHost("server2"))

	hosts, err := GetAllHost(db)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(hosts) != 2 {
		t.Fatalf("expected 2 hosts, got %d", len(hosts))
	}
}

func TestGetAllHost_FieldMapping(t *testing.T) {
	db := setupTestDB(t)
	h := model.SSHHost{
		Alias:    "web",
		Host:     "10.0.0.1",
		User:     "admin",
		Port:     2222,
		AuthType: "Public Key",
		KeyPath:  "~/.ssh/id_rsa",
	}
	_ = AddHost(db, h)

	hosts, _ := GetAllHost(db)
	got := hosts[0]

	if got.Alias != h.Alias || got.Host != h.Host || got.User != h.User ||
		got.Port != h.Port || got.AuthType != h.AuthType || got.KeyPath != h.KeyPath {
		t.Fatalf("field mismatch:\n  expected: %+v\n  got:      %+v", h, got)
	}
}

// --- DeleteHost ---

func TestDeleteHost_Success(t *testing.T) {
	db := setupTestDB(t)
	_ = AddHost(db, newTestHost("todelete"))

	hosts, _ := GetAllHost(db)
	err := DeleteHost(db, hosts[0].ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	remaining, _ := GetAllHost(db)
	if len(remaining) != 0 {
		t.Fatalf("expected 0 hosts after delete, got %d", len(remaining))
	}
}

func TestDeleteHost_NotFound(t *testing.T) {
	db := setupTestDB(t)

	err := DeleteHost(db, 9999)
	if err == nil {
		t.Fatal("expected error for non-existent ID, got nil")
	}
}

// --- UpdateHost ---

func TestUpdateHost_Success(t *testing.T) {
	db := setupTestDB(t)
	_ = AddHost(db, newTestHost("original"))

	hosts, _ := GetAllHost(db)
	h := hosts[0]
	h.Alias = "updated"
	h.Port = 3333

	err := UpdateHost(db, h)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	after, _ := GetAllHost(db)
	if after[0].Alias != "updated" || after[0].Port != 3333 {
		t.Fatalf("update not applied: %+v", after[0])
	}
}

func TestUpdateHost_DuplicateAlias(t *testing.T) {
	db := setupTestDB(t)
	_ = AddHost(db, newTestHost("a"))
	_ = AddHost(db, newTestHost("b"))

	hosts, _ := GetAllHost(db)
	h := hosts[1]
	h.Alias = "a" // alias "a" already exists

	err := UpdateHost(db, h)
	if err == nil {
		t.Fatal("expected error for duplicate alias on update, got nil")
	}
}

// --- ResetDB ---

func TestResetDB_ClearsAllData(t *testing.T) {
	db := setupTestDB(t)
	_ = AddHost(db, newTestHost("s1"))
	_ = AddHost(db, newTestHost("s2"))
	_ = AddHost(db, newTestHost("s3"))

	err := ResetDB(db)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	hosts, _ := GetAllHost(db)
	if len(hosts) != 0 {
		t.Fatalf("expected 0 hosts after reset, got %d", len(hosts))
	}
}
