package ui

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"ok2pus/internal/db"
	"ok2pus/internal/model"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func AddNewHostInteractive(d *sql.DB) {
	var err error

	validate := func(input string) error {
		if len(strings.TrimSpace(input)) < 1 {
			return fmt.Errorf("this field is required")
		}
		return nil
	}

	alias := promptNewAlias(d)
	if alias == "" {
		return
	}

	promptAuth := promptui.Select{
		Label: "Select ssh authentication method",
		Items: []string{"[1] Password", "[2] Public Key", "[3] Back"},
	}

	fmt.Println()
	_, authResult, authErr := promptAuth.Run()
	if authErr != nil || authResult == "[3] Back" {
		return
	}

	var path string
	var auth string

	switch authResult {
	case "[1] Password":
		auth = "Password"
	case "[2] Public Key":
		auth = "Public Key"
		fmt.Println()
		promptPath := promptui.Prompt{Label: "Key file path", Default: "~/.ssh/id_rsa", Validate: validate}
		path, err = promptPath.Run()
		if err != nil {
			return
		}
	}

	fmt.Println()
	promptUser := promptui.Prompt{Label: "Username", Validate: validate}
	user, err := promptUser.Run()
	if err != nil {
		return
	}

	fmt.Println()
	promptHost := promptui.Prompt{Label: "Host Address", Validate: validate}
	host, err := promptHost.Run()
	if err != nil {
		return
	}

	validatePort := func(input string) error {
		port, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			return fmt.Errorf("port must be a number")
		}
		if len(strings.TrimSpace(input)) < 1 {
			return fmt.Errorf("this field is required")
		}
		if port < 1 || port > 65535 {
			return fmt.Errorf("port must be between 1 and 65535")
		}
		return nil
	}

	fmt.Println()
	promptPort := promptui.Prompt{Label: "Port", Default: "22", Validate: validatePort}

	portStr, err := promptPort.Run()
	if err != nil {
		return
	}

	port, _ := strconv.Atoi(portStr)

	err = db.AddHost(d, model.SSHHost{Alias: alias, Host: host, User: user, Port: port, AuthType: auth, KeyPath: path})
	if err != nil {
		color.Red("Error: %v\n", err)
	} else {
		color.New(color.FgGreen, color.Bold).Println("\nSuccessfully added!")
	}
}

func promptNewAlias(d *sql.DB) string {
	validate := func(input string) error {
		if len(input) < 1 {
			return nil
		}
		if db.ExistsAlias(d, input) {
			return errors.New("The alias already exists. Please choose another one.")
		}
		return nil
	}

	for {
		prompt := promptui.Prompt{
			Label:    "Alias(title)",
			Validate: validate,
		}

		result, err := prompt.Run()
		if err != nil {
			return ""
		}
		if len(result) < 1 {
			continue
		}
		return result
	}
}
