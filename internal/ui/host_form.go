package ui

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"ok2pus/internal/db"
	"ok2pus/internal/model"

	"github.com/manifoldco/promptui"
)

func AddNewHostInteractive(d *sql.DB) {
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

	validatePort := func(input string) error {
		if len(strings.TrimSpace(input)) < 1 {
			return fmt.Errorf("this field is required")
		}
		port, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			return fmt.Errorf("port must be a number")
		}
		if port < 1 || port > 65535 {
			return fmt.Errorf("port must be between 1 and 65535")
		}
		return nil
	}

	promptPort := promptui.Prompt{Label: "Port", Default: "22", Validate: validatePort}
	portStr, _ := promptPort.Run()
	port, _ := strconv.Atoi(portStr)

	if alias == "" || host == "" {
		return
	}

	err := db.AddHost(d, model.SSHHost{Alias: alias, Host: host, User: user, Port: port, AuthType: auth, KeyPath: path})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Successfully Added!")
	}
}
