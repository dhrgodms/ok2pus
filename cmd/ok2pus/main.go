package main

import (
	"fmt"
	"os"

	"ok2pus/internal/db"
	"ok2pus/internal/ui"

	"github.com/manifoldco/promptui"
)

func main() {
	d := db.InitDB()
	defer d.Close()
	ui.ShowLogo()

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
			ui.ShowHostListMenu(d)
		case "[2] Add New Host":
			ui.AddNewHostInteractive(d)
		case "[3] Options":
			ui.ShowOptionsMenu(d)
		case "[q] Quit":
			fmt.Println("\nGoodbye!")
			return
		}
	}
}
