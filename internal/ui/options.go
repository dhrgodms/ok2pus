package ui

import (
	"database/sql"
	"fmt"

	"ok2pus/internal/db"

	"github.com/manifoldco/promptui"
)

func ShowOptionsMenu(d *sql.DB) {
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
		resetDatabase(d)
	case "[2] Drop Database":
		dropDatabase(d)
	}
}

func resetDatabase(d *sql.DB) {
	confirmPrompt := promptui.Prompt{
		Label:     "Are you sure you want to RESET all data?",
		IsConfirm: true,
	}

	_, err := confirmPrompt.Run()
	if err != nil {
		fmt.Println("Reset cancelled.")
		return
	}

	err = db.ResetDB(d)
	if err != nil {
		fmt.Printf("Error during reset: %v\n", err)
	} else {
		fmt.Println("Successfully reset database. All host information has been cleared.")
	}
}

func dropDatabase(d *sql.DB) {
	confirmPrompt := promptui.Prompt{
		Label:     "Are you sure you want to DROP the database?",
		IsConfirm: true,
	}

	_, err := confirmPrompt.Run()
	if err != nil {
		fmt.Println("Drop cancelled.")
		return
	}

	db.DropDB(d)
}
