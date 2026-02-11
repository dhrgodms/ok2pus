package main

import (
	"database/sql"
	"fmt"
)

type SSHHost struct {
	ID int
	Alias string
	Host string
	User string
	Port int
}

func addHost(db *sql.DB, h SSHHost) error {
	query := `INSERT INTO ssh_hosts (alias, host, user, port) VALUES (?, ?, ?, ?);`
	_, err := db.Exec(query, h.Alias, h.Host, h.User, h.Port)
	return err
}

func getAllHost(db *sql.DB) ([]SSHHost, error) {
	rows, err := db.Query("SELECT id, alias, host, user, port FROM ssh_hosts")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var hosts []SSHHost
	for rows.Next() {
		var h SSHHost
		if err := rows.Scan(&h.ID, &h.Alias, &h.Host, &h.User, &h.Port); err != nil {
			return nil, err
		}
		hosts = append(hosts, h)
	}
	return hosts, nil
}

func deleteHost(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM ssh_hosts WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no host found with ID %d", id)
	}
	return nil
}

func updateHost() {}