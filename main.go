package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"

	"github.com/manifoldco/promptui"
)

func main() {
	db := initDB()
	defer db.Close()
	showLogo()

	for {
		prompt := promptui.Select{
			Label:        "Select",
			Items:        []string{"[1] List Hosts", "[2] Add New Host", "[3] Options", "[q] Quit"},
			HideSelected: true,
		}

		_, result, err := prompt.Run()
		if err != nil {
			if err == promptui.ErrInterrupt {
				fmt.Println("\nStopped by user.")
				os.Exit(0)
			}
			continue
		}

		switch result {
		case "[1] List Hosts":
			showHostListMenu(db)
		case "[2] Add New Host":
			addNewHostInteractive(db)
		case "[3] Options":
			showOptionsMenu(db)
		case "[q] Quit":
			fmt.Println("\nGoodbye!")
			return
		}
	}
}

func showOptionsMenu(db *sql.DB) {
	prompt := promptui.Select{
		Label:        "Select Options",
		Items:        []string{"[1] Reset Database", "[2] Drop Database", "[3] Back"},
		HideSelected: true,
	}

	_, result, err := prompt.Run()

	if err != nil || result == "[3] Back" {
		return
	}

	switch result {
	case "[1] Reset Database":
		resetDatabase(db)
	case "[2] Drop Database":
		dropDatabase(db)
	}
}

func addNewHostInteractive(db *sql.DB) {
	validate := func(input string) error {
		if len(strings.TrimSpace(input)) < 1 {
			return fmt.Errorf("this field is required")
		}
		return nil
	}

	promptAuth := promptui.Select{
		Label:        "Select ssh authentication method",
		Items:        []string{"[1] Password", "[2] Public Key", "[3] Back"},
		HideSelected: true,
	}

	_, authResult, authErr := promptAuth.Run()
	if authErr != nil || authResult == "[3] Back" {
		return
	}

	promptAlias := promptui.Prompt{Label: "Alias", Validate: validate}
	alias, _ := promptAlias.Run()

	var path string
	var auth string

	switch authResult {
	case "[1] Password":
		auth = "Password"
	case "[2] Public Key":
		auth = "Public Key"
		promptPath := promptui.Prompt{Label: "Key file path", Default: "~/.ssh/id_rsa", Validate: validate}
		path, _ = promptPath.Run()
	}

	promptUser := promptui.Prompt{Label: "Username", Validate: validate}
	user, _ := promptUser.Run()

	promptHost := promptui.Prompt{Label: "Host Address", Validate: validate}
	host, _ := promptHost.Run()

	promptPort := promptui.Prompt{Label: "Port", Default: "22", Validate: validate}
	portStr, _ := promptPort.Run()
	port, _ := strconv.Atoi(portStr)

	if alias == "" || host == "" {
		return
	}

	err := addHost(db, SSHHost{Alias: alias, Host: host, User: user, Port: port, AuthType: auth, KeyPath: path})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Successfully Added!")
	}
}

func showHostListMenu(db *sql.DB) {
	hosts, _ := getAllHost(db)
	if len(hosts) == 0 {
		fmt.Println("No hosts found.")
		return
	}

	var items []string
	for _, h := range hosts {
		items = append(items, fmt.Sprintf("[%s] %s@%s:%d (%s)", h.Alias, h.User, h.Host, h.Port, h.AuthType))
	}
	items = append(items, "Back")

	prompt := promptui.Select{
		Label:        "Select a Host",
		Items:        items,
		Size:         5,
		HideSelected: true,
	}

	index, result, err := prompt.Run()
	if err != nil || result == "Back" {
		return
	}

	selectedHost := hosts[index]
	showActionMenu(db, selectedHost)
}

func showActionMenu(db *sql.DB, host SSHHost) {
	prompt := promptui.Select{
		Label:        fmt.Sprintf("Action for [%s]", host.Alias),
		Items:        []string{"[1] connect", "[2] edit", "[3] delete", "[4] Back"},
		HideSelected: true,
	}

	_, result, _ := prompt.Run()

	switch result {
	case "[1] connect":
		connectHost(host)
	case "[2] edit":
		openEditor(db, host)
	case "[3] delete":
		confirmPrompt := promptui.Prompt{
			Label:     fmt.Sprintf("Are you sure you want to delete [%s]?", host.Alias),
			IsConfirm: true,
		}
		_, err := confirmPrompt.Run()
		if err != nil {
			return // 취소 시 복귀
		}

		err = deleteHost(db, host.ID)
		if err != nil {
			fmt.Printf("Delete failed: %v\n", err)
		} else {
			fmt.Println("Successfully deleted.")
		}
	case "[4] Back":
		return
	}
}

type editFinishedMsg struct {
	content string
	err     error
}

func openEditor(db *sql.DB, h SSHHost) {
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

	err = updateHost(db, updated)
	if err != nil {
		fmt.Println("Error occurred during updating the cofig. err:", err)
	} else {
		fmt.Println("Successfully updating the cofig.")
	}
}

func parseHostConfig(content string, original SSHHost) (SSHHost, error) {
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
