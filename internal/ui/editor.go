package ui

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"ok2pus/internal/db"
	"ok2pus/internal/model"
)

func OpenEditor(d *sql.DB, h model.SSHHost) {
	currentConfig := fmt.Sprintf("alias=%s\nhost=%s\nuser=%s\nport=%d\nauth_type(Password/Public Key)=%s\nkey_path=%s", h.Alias, h.Host, h.User, h.Port, h.AuthType, h.KeyPath)

	tempFile, err := os.CreateTemp("", "host-config-*.txt")
	if err != nil {
		fmt.Println("Failed to create a temp file. err:", err)
		return
	}
	defer os.Remove(tempFile.Name())
	tempFile.WriteString(currentConfig)
	tempFile.Close()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error occurred during executing the editor. err:", err)
		return
	}

	updatedConfig, _ := os.ReadFile(tempFile.Name())

	updated, err := parseHostConfig(string(updatedConfig), h)
	if err != nil {
		fmt.Print("Error occurred during parsing. err:", err)
		return
	}

	err = db.UpdateHost(d, updated)
	if err != nil {
		fmt.Println("Error occurred during updating the cofig. err:", err)
	} else {
		fmt.Println("Successfully updating the cofig.")
	}
}

func parseHostConfig(content string, original model.SSHHost) (model.SSHHost, error) {
	trimmedContent := strings.TrimSpace(strings.TrimSpace(content))
	if trimmedContent == "" {
		return original, fmt.Errorf("There's no change.")
	}

	updated := original
	lines := strings.Split(trimmedContent, "\n")
	foundFields := make(map[string]bool)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, "=") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		val := strings.TrimSpace(parts[1])

		switch key {
		case "alias":
			updated.Alias = val
			if val != "" {
				foundFields["alias"] = true
			}
		case "user":
			updated.User = val
			if val != "" {
				foundFields["user"] = true
			}
		case "host":
			updated.Host = val
			if val != "" {
				foundFields["host"] = true
			}
		case "port":
			port, err := strconv.Atoi(val)
			if err == nil {
				if port < 1 || port > 65535 {
					return original, fmt.Errorf("port must be between 1 and 65535, got %d", port)
				}
				updated.Port = port
			}
		case "auth_type":
			updated.AuthType = val
		case "key_path":
			updated.KeyPath = val
		}
	}

	if !foundFields["alias"] || !foundFields["host"] {
		return original, fmt.Errorf("Alias & Host fields cannot be empty.")
	}
	return updated, nil
}
