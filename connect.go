package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func connectHost(h SSHHost) {
	dest := fmt.Sprintf("%s@%s", h.User, h.Host)

	cmd := exec.Command("ssh", dest, "-p", strconv.Itoa(h.Port))

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("\n--- Connecting to [%s] (%s) ---\n", h.Alias, h.Host)

	err := cmd.Run()
	if err != nil {
		fmt.Printf("\n[SSH Error] %v\n", err)
	}

	fmt.Printf("\n--- Connection to [%s] closed ---\n", h.Alias)
}