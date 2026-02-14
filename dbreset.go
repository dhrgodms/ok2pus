package main

import (
	"database/sql"
	"fmt"

	"github.com/manifoldco/promptui"
)

func resetDatabase(db *sql.DB) {
	confirmPrompt := promptui.Prompt{
		Label:     "Are you sure you want to RESET all data?",
		IsConfirm: true,
	}

	_, err := confirmPrompt.Run()
	if err != nil {
		fmt.Println("Reset cancelled.")
		return
	}

	err = resetDB(db)
	if err != nil {
		fmt.Printf("Error during reset: %v\n", err)
	} else {
		fmt.Println("Successfully reset database. All host information has been cleared.")
	}
}

func dropDatabase(db *sql.DB) {
	confirmPrompt := promptui.Prompt{
		Label:     "Are you sure you want to DROP the database?",
		IsConfirm: true,
	}

	_, err := confirmPrompt.Run()
	if err != nil {
		fmt.Println("Drop cancelled.")
		return
	}

	dropDB(db)
}
