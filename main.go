package main

import (
	"database/sql"
	"fmt"
	"os"
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
			Items:        []string{"1. List Hosts", "2. Add New Host", "q. Quit"},
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
		case "1. List Hosts":
			showHostListMenu(db)
		case "2. Add New Host":
			addNewHostInteractive(db)
		case "q. Quit":
			fmt.Println("\nGoodbye!")
			return
		}
	}
}

func addNewHostInteractive(db *sql.DB) {
	validate := func(input string) error {
		if len(strings.TrimSpace(input)) < 1 {
			return fmt.Errorf("this field is required")
		}
		return nil
	}

	promptAlias := promptui.Prompt{Label: "Alias", Validate: validate}
	alias, _ := promptAlias.Run()

	promptHost := promptui.Prompt{Label: "Host Address", Validate: validate}
	host, _ := promptHost.Run()

	promptUser := promptui.Prompt{Label: "Username", Validate: validate}
	user, _ := promptUser.Run()

	promptPort := promptui.Prompt{Label: "Port", Default: "22", Validate: validate}
	portStr, _ := promptPort.Run()
	port, _ := strconv.Atoi(portStr)

	if alias == "" || host == "" {
		return
	}

	err := addHost(db, SSHHost{Alias: alias, Host: host, User: user, Port: port})
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
		items = append(items, fmt.Sprintf("[%s] %s@%s:%d", h.Alias, h.User, h.Host, h.Port))
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
		Items:        []string{"connect", "edit", "delete", "Back"},
		HideSelected: true,
	}

	_, result, _ := prompt.Run()

	switch result {
	case "connect":
		connectHost(host)
	case "edit":
		// TODO: implement
		updateHost(host)
		fmt.Println("Updated.")
	case "delete":
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
	}
}
