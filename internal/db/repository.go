package db

import (
	"database/sql"
	"fmt"

	"ok2pus/internal/model"
)

func AddHost(db *sql.DB, h model.SSHHost) error {
	query := `INSERT INTO ssh_hosts (alias, host, user, port, auth_type, key_path) VALUES (?, ?, ?, ?, ?, ?);`
	_, err := db.Exec(query, h.Alias, h.Host, h.User, h.Port, h.AuthType, h.KeyPath)
	return err
}

func GetAllHost(db *sql.DB) ([]model.SSHHost, error) {
	rows, err := db.Query("SELECT id, alias, host, user, port, auth_type, key_path FROM ssh_hosts")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var hosts []model.SSHHost
	for rows.Next() {
		var h model.SSHHost
		if err := rows.Scan(&h.ID, &h.Alias, &h.Host, &h.User, &h.Port, &h.AuthType, &h.KeyPath); err != nil {
			return nil, err
		}
		hosts = append(hosts, h)
	}
	return hosts, nil
}

func DeleteHost(db *sql.DB, id int) error {
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

func UpdateHost(db *sql.DB, h model.SSHHost) error {
	query := `
	UPDATE ssh_hosts
	SET alias=?, host=?, user=?, port=?, auth_type=?, key_path=?
	WHERE id=?;`

	_, err := db.Exec(query, h.Alias, h.Host, h.User, h.Port, h.AuthType, h.KeyPath, h.ID)
	return err
}

func ResetDB(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM ssh_hosts; DELETE FROM sqlite_sequence WHERE name='ssh_hosts';")
	return err
}

func ExistsAlias(d *sql.DB, alias string) bool {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM ssh_hosts WHERE alias = ?)"
	err := d.QueryRow(query, alias).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}
