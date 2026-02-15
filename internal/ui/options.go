package ui

import (
	"database/sql"

	"ok2pus/internal/db"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func ShowOptionsMenu(d *sql.DB) {
	prompt := promptui.Select{
		Label: "Select Options",
		Items: []string{"[1] Reset Database", "[2] Drop Database", "[3] Back"},
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
		color.Red("Reset cancelled.")
		return
	}

	err = db.ResetDB(d)
	if err != nil {
		color.Red("Error during reset: %v\n", err)
	} else {
		color.New(color.Bold, color.FgGreen).Print("Successfully reset database. ")
		color.White("All host information has been cleared.\n")
	}
}

func dropDatabase(d *sql.DB) {
	confirmPrompt := promptui.Prompt{
		Label:     "Are you sure you want to DROP the database?(remove db file)",
		IsConfirm: true,
	}

	_, err := confirmPrompt.Run()
	if err != nil {
		color.Red("Drop cancelled.")
		return
	}

	db.DropDB(d)
}
