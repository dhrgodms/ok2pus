package ui

import (
	"database/sql"
	"fmt"

	"ok2pus/internal/db"
	"ok2pus/internal/model"
	"ok2pus/internal/ssh"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func ShowHostListMenu(d *sql.DB) {
	hosts, _ := db.GetAllHost(d)
	if len(hosts) == 0 {
		color.Yellow("No hosts found.")
		return
	}

	var items []string
	for _, h := range hosts {
		items = append(items, fmt.Sprintf("[%s] %s@%s:%d (%s)", color.New(color.Bold, color.FgCyan).Sprint(h.Alias), h.User, h.Host, h.Port, h.AuthType))
	}
	items = append(items, "Back")

	prompt := promptui.Select{
		Label: "Select a Host",
		Items: items,
		Size:  5,
	}

	index, result, err := prompt.Run()
	if err != nil || result == "Back" {
		return
	}

	selectedHost := hosts[index]
	showActionMenu(d, selectedHost)
}

func showActionMenu(d *sql.DB, host model.SSHHost) {
	prompt := promptui.Select{
		Label: fmt.Sprintf("Action for [%s]", host.Alias),
		Items: []string{"[1] connect", "[2] edit", "[3] delete", "[4] Back"},
	}

	_, result, _ := prompt.Run()

	switch result {
	case "[1] connect":
		ssh.ConnectHost(host)
	case "[2] edit":
		OpenEditor(d, host)
	case "[3] delete":
		deleteHost(d, host)
	case "[4] Back":
		return
	}
}

func deleteHost(d *sql.DB, host model.SSHHost) {
	fmt.Println()
	confirmPrompt := promptui.Prompt{
		Label:     fmt.Sprintf("Are you sure you want to delete [%s]", host.Alias),
		IsConfirm: true,
	}
	_, err := confirmPrompt.Run()
	if err != nil {
		return
	}

	err = db.DeleteHost(d, host.ID)
	if err != nil {
		color.Red("Delete failed: %v\n", err)
		return
	}
	color.New(color.Bold, color.FgGreen).Println("\nSuccessfully deleted.")
}
